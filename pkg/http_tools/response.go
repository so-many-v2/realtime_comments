package http_tools

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, code int, err error) {
	JSON(w, code, map[string]string{"error": err.Error()})
}
