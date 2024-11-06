package geonode

type Proxy struct {
	Ip           string  `json:"ip"`
	Port         string  `json:"port"`
	Uptime       float32 `json:"uptime"`
	ResponseTime float32 `json:"responseTime"`
}
