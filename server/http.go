package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func httpError(w http.ResponseWriter, statusCode int, err error) {
	log.Printf("%s: %s\n", http.StatusText(statusCode), err)
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func httpOkJson(w http.ResponseWriter, payload any) (err error) {
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(payload)
}
