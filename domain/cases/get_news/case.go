package get_news

import (
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/antchfx/htmlquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/bot"
)

var (
	site        = "https://kakoysegodnyaprazdnik.ru/"
	headerName  = "User-Agent"
	header      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
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
	bot.Bot.Send(tgbotapi.NewMessage(chatId, "Загружаю новости (примерно 10 секунд)..."))
	openFirefox()

	req, err := http.NewRequest("GET", site, nil)
	if err != nil {
		log.Println("Error while get to " + site + " " + err.Error())
	}

	req.Header.Add(headerName, header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error while do request " + site + " " + err.Error())
	}
	defer resp.Body.Close()

	log.Println("Request to "+site+" processed with code ", resp.StatusCode)

	node, err := htmlquery.Parse(resp.Body)
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

func openFirefox() {
	cmd := exec.Command("firefox", "--headless", "https://kakoysegodnyaprazdnik.ru/")

	go func() {
		err := cmd.Run()
		if err != nil {
			log.Println("Error while open firefox " + err.Error())
		}
	}()

	time.Sleep(60 * time.Second)
	err := cmd.Process.Kill()
	if err != nil {
		log.Println("Error while kill firefox " + err.Error())
	}
}
