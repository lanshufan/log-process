package types

type LogFormat struct {
	Ip               string  `json:"ip"`
	Date             string  `json:"date"`
	Method           string  `json:"method"`
	RequestPath      string  `json:"requestPath"`
	RequestSize      int     `json:"requestSize"`
	UpstreamAddr     string  `json:"upstreamAddr"`
	UpstreamTime     float64 `json:"upstreamTime"`
	ResponseTime     float64 `json:"responseTime"`
	ResponseStatus   int     `json:"responseStatus"`
	ResponseBodySize int     `json:"responseBodySize"`
}
