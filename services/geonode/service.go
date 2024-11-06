package geonode

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/paavill/awesome-tagger-bot/domain/models"
	"github.com/paavill/awesome-tagger-bot/domain/services"
	"github.com/patrickmn/go-cache"
)

func New(host string) services.GetProxy {
	return &geonode{
		host:  host,
		cache: cache.New(24*time.Hour, 0),
	}
}

type geonode struct {
	host  string //https://proxylist.geonode.com
	cache *cache.Cache
}

func (g *geonode) GetProxyList() ([]*models.Proxy, error) {
	path := "/api/proxy-list?protocols=socks5&limit=50&page=1&sort_by=upTime&sort_type=desc"
	url, err := url.Parse(fmt.Sprintf("%s%s", g.host, path))
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    url,
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	body := struct{
		Data []Proxy `json:"data"`
	}{
		Data: []Proxy{},
	}

	err = json.Unmarshal(bodyRaw, &body)
	if err != nil {
		return nil, err
	}

	result := make([]*models.Proxy, len(body.Data))
	for i := 0; i < len(body.Data); i++ {
		result[i] = &models.Proxy{
			Ip:           body.Data[i].Ip,
			Port:         body.Data[i].Port,
			Uptime:       body.Data[i].Uptime,
			ResponseTime: body.Data[i].ResponseTime,
		}
	}

	return result, nil
}

func (g *geonode) GetProxyListCached() ([]*models.Proxy, error) {
	rawList, exists := g.cache.Get("proxy_list")
	if exists {
		list := rawList.([]*models.Proxy)
		return list, nil
	}
	list, err := g.GetProxyList()
	if err != nil {
		return nil, err
	}
	g.cache.Set("proxy_list", list, cache.DefaultExpiration)
	return list, nil
}
