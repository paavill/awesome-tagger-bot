package get_news

import (
	"log"
	"net/http"
	"time"

	"github.com/antchfx/htmlquery"
)

var (
	site       = "https://kakoysegodnyaprazdnik.ru/"
	headerName = "User-Agent"
	header     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
	h1         = struct {
		k string
		v string
	}{
		"authority",
		"kakoysegodnyaprazdnik.ru",
	}
	h2 = struct {
		k string
		v string
	}{
		"method",
		"GET",
	}
	h3 = struct {
		k string
		v string
	}{
		"path",
		"/",
	}
	h4 = struct {
		k string
		v string
	}{
		"scheme",
		"https",
	}
	cachedTitle = ""
	cachedNews  = []string{}
	cachedDay   = -1
)

func Run() (string, []string, error) {
	t, n, ok := getCached()
	if ok {
		log.Println("Get cached news")
		return t, n, nil
	}

	req, err := http.NewRequest("GET", site, nil)
	if err != nil {
		log.Println("Error while get to " + site + " " + err.Error())
	}

	req.Header.Add(headerName, header)
	req.Header.Add(h1.k, h1.v)
	req.Header.Add(h2.k, h2.v)
	req.Header.Add(h3.k, h3.v)
	req.Header.Add(h4.k, h4.v)

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
