package get_news

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/antchfx/htmlquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/paavill/awesome-tagger-bot/domain/context"
	"golang.org/x/net/html"
)

var (
	site                    = "https://kakoysegodnyaprazdnik.ru/"
	cachedTitle             = ""
	cachedNews              = []string{}
	cachedDay               = -1
	mux         *sync.Mutex = &sync.Mutex{}
	muxLocked               = false
)

func ClearCache() {
	cachedTitle = ""
	cachedNews = []string{}
	cachedDay = -1
}

func Run(ctx context.Context, chatId int64) (string, []string, error) {
	var body string = ""
	defer func() {
		if r := recover(); r != nil {
			fileUuid := uuid.New().String()
			os.WriteFile("./"+fileUuid+".html", []byte(body), 0777)
			log.Println("Recovered in f", r)
		}
	}()

	//if muxLocked {
	//bot.Bot.Send(tgbotapi.NewMessage(chatId, "Уже загружаю, осталось чуть-чуть"))
	//}

	mux.Lock()
	muxLocked = true
	defer func() {
		mux.Unlock()
		muxLocked = false
	}()

	t, n, ok := getCached()
	if ok {
		ctx.Logger().Info("Get cached news")
		return t, n, nil
	}

	ctx.Logger().Info("Get news from " + site)
	ctx.Services().Bot().Send(tgbotapi.NewMessage(chatId, "Загружаю новости (примерно 30 секунд)..."))

	body = getHtml()
	bodyReader := strings.NewReader(body)

	node, err := htmlquery.Parse(bodyReader)
	if err != nil {
		ctx.Logger().Error("Error while parse html " + site + " " + err.Error())
	}

	title, err := htmlquery.Query(node, "//html//body//div[1]//h2")
	if err != nil {
		ctx.Logger().Error("Error while get title " + site + " " + err.Error())
	}
	titleText := htmlquery.InnerText(title)

	news := []string{}
	rn, err := htmlquery.Query(node, "/html/body/div[1]/div[1]/div")
	if err != nil {
		ctx.Logger().Error("Error while get news " + site + " " + err.Error())
	}

	newsNodes := findNodesWithAttrValue(rn, "itemprop", "suggestedAnswer", "acceptedAnswer")

	mainNodes := []*html.Node{}
	for _, node := range newsNodes {
		mainNodes = append(mainNodes, findNodesWithAttrValue(node, "class", "main")...)
	}

	for _, node := range mainNodes {
		textNodes := findNodesWithAttrValue(node, "itemprop", "text")
		superTextNodes := findNodesWithAttrValue(node, "class", "super")
		hrefTextNodes := findNodesWithAttrValue(node, "class", "prazdnik_info")
		newText := ""
		if len(textNodes) > 0 {
			newText += "• " + htmlquery.InnerText(textNodes[0])
		} else if len(hrefTextNodes) > 0 {
			hrefNode := hrefTextNodes[0]
			href := htmlquery.SelectAttr(hrefNode, "href")
			href = "https://kakoysegodnyaprazdnik.ru" + href
			if href != "" {
				newText += "• " + fmt.Sprintf("<a href=\"%s\"><i><b>%s</b></i></a>", href, htmlquery.InnerText(hrefNode))
			} else {
				newText += "• " + htmlquery.InnerText(hrefNode)
			}
		}
		if len(superTextNodes) > 0 {
			newText += " " + htmlquery.InnerText(superTextNodes[0])
		}
		if newText != "" {
			news = append(news, newText)
		}
	}

	setCached(titleText, news)

	return titleText, news, err
}

func findNodesWithAttrValue(node *html.Node, attrName string, attrValue ...string) []*html.Node {
	sibling := node.FirstChild
	result := []*html.Node{}
	for sibling != nil {
		attrs := map[string]string{}
		for _, v := range sibling.Attr {
			attrs[v.Key] = v.Val
		}

		if v, ok := attrs[attrName]; ok  {
			vOk := false
			for _, attrV := range attrValue {
				vOk = vOk || v == attrV
				if vOk {
					break
				}
			}
			if vOk {
				result = append(result, sibling)
			}
		}
		sibling = sibling.NextSibling
	}
	return result
}

func getCached() (string, []string, bool) {
	n := time.Now()

	if cachedDay == n.Day() {
		return cachedTitle, cachedNews, true
	}

	return "", nil, false
}

func setCached(title string, news []string) {
	cachedTitle = title
	cachedNews = news

	n := time.Now()
	cachedDay = n.Day()
}

func getHtml() string {
	cmd := exec.Command("python3", "get_news.py")

	output, err := cmd.Output()
	if err != nil {
		log.Println("Error while open firefox " + err.Error())
	}

	return string(output)
}
