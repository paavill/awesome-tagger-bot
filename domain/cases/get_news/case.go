package get_news

import (
	"log"
	"net/http"

	"github.com/antchfx/htmlquery"
)

var (
	site       = "https://kakoysegodnyaprazdnik.ru/"
	headerName = "User-Agent"
	header     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
)

func Run() (string, []string, error) {
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

	return titleText, news, err
}
