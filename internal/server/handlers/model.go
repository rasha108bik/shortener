package handlers

type ReqCreateShorten struct {
	URL string `json:"url"`
}

type RespReqCreateShorten struct {
	Result string `json:"result"`
}
