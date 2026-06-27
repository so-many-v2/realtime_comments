package http_tools

import (
	"encoding/json"
	"net/http"
)

func NewDecoder(w http.ResponseWriter, req *http.Request) *json.Decoder {
	dec := json.NewDecoder(http.MaxBytesReader(w, req.Body, MaxBodySize))
	return dec
}
