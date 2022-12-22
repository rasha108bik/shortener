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
