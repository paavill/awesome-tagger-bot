package get_news

import (
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/bot"
)

var (
	site        = "https://kakoysegodnyaprazdnik.ru/"
	cachedTitle = ""
	cachedNews  = []string{}
	cachedDay   = -1
)

func Run(chatId int64) (string, []string, error) {
	t, n, ok := getCached()
	if ok {
		log.Println("Get cached news")
		return t, n, nil
	}

	log.Println("Get news from " + site)
	bot.Bot.Send(tgbotapi.NewMessage(chatId, "Загружаю новости (примерно 30 секунд)..."))

	body := getHtml()
	bodyReader := strings.NewReader(body)

	node, err := htmlquery.Parse(bodyReader)
	if err != nil {
		log.Println("Error while parse html " + site + " " + err.Error())
	}

	title, err := htmlquery.Query(node, "//html//body//div[1]//h2")
	if err != nil {
		log.Println("Error while get title " + site + " " + err.Error())
	}
	titleText := htmlquery.InnerText(title)

	news := []string{}
	rn, err := htmlquery.Query(node, "/html/body/div[1]/div[1]/div")
	if err != nil {
		log.Println("Error while get news " + site + " " + err.Error())
	}

	sibbling := rn.FirstChild
	for sibbling.NextSibling != nil {
		attrs := map[string]string{}
		for _, v := range sibbling.Attr {
			attrs[v.Key] = v.Val
		}
		if v, ok := attrs["itemprop"]; ok && v == "suggestedAnswer" {
			text := htmlquery.InnerText(sibbling)
			news = append(news, text)
		}
		sibbling = sibbling.NextSibling
	}

	setCached(titleText, news)

	return titleText, news, err
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
