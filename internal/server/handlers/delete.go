package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rasha108bik/tiny_url/internal/utility"
)

// DeleteURLs delete urls which include in reqBody.
func (h *handler) DeleteURLs(w http.ResponseWriter, r *http.Request) {
	var reqBody []string
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		h.log.Err(err).Msg("decode body failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.log.Info().Msgf("reqBody: %#v", reqBody)

	inputCh := make(chan string)
	go func(reqBody []string) {
		for _, v := range reqBody {
			inputCh <- v
		}

		close(inputCh)
	}(reqBody)

	fanOuts := utility.FanOut(inputCh, len(reqBody))
	for _, fanOutCh := range fanOuts {
		go func(fanOutCh chan string) {
			ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
			defer cancel()

			err = h.storage.DeleteURLByShortURL(ctx, <-fanOutCh)
			if err != nil {
				h.log.Err(err).Msg("DeleteURLsbyShortURL failed")
			}
		}(fanOutCh)
	}

	w.WriteHeader(http.StatusAccepted)
}
