package handlers

/* HTTP request utils */

import (
	"encoding/json"
	"net/http"
)

func responseJSON(w http.ResponseWriter, payload interface{}) {
	// Convert struct to JSON sending the response
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func responseError(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusInternalServerError)
}

func responseText(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}
