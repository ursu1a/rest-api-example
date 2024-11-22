package handlers
/* HTTP request utils */

import (
	"net/http"
	"encoding/json"
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}