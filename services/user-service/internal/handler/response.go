package handler

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response{
		Success: status < 400,
		Data:    data,
	})
}
