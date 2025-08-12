package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code" example:"401"`
	Message string `json:"error" example:"Inavalid credentials"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {

	RespondWithJSON(w, code, ErrorResponse{
		Code:    code,
		Message: msg,
	})
}
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal Server Error"}`))
		return
	}
	w.WriteHeader(code)
	if _, err := w.Write(dat); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
