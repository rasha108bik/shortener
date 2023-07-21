package models

// models for api requests
type (
	// ReqCreateShorten model for CreateShorten request.
	ReqCreateShorten struct {
		URL string `json:"url"`
	}

	// RespReqCreateShorten model for response CreateShorten request.
	RespReqCreateShorten struct {
		Result string `json:"result"`
	}

	// RespGetOriginalURLs model for FetchURLs request.
	RespGetOriginalURLs struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	// ReqShortenBatch model for ShortenBatch request.
	ReqShortenBatch struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	// RespShortenBatch model for response ShortenBatch request.
	RespShortenBatch struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}
)
