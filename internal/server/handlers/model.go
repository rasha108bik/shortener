package handlers

type ReqCreateShorten struct {
	URL string `json:"url"`
}

type RespReqCreateShorten struct {
	Result string `json:"result"`
}

type RespGetOriginalURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ReqShortenBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type RespShortenBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
