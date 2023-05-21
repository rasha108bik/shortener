package handlers

import "net/http"

// ErrorHandler method which not include API methods
func (h *handler) ErrorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "wrong method", http.StatusBadRequest)
}
