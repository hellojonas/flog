package api

import (
	"encoding/json"
	"net/http"
)

type HttpMessageResponse struct {
	Message string `json:"message"`
}

func sendResponse(w http.ResponseWriter, status int) error {
	return sendJson(w, status, nil)
}

func sendJson(w http.ResponseWriter, status int, body interface{}) error {
	w.Header().Add("Content-type", "application/json")

	var data []byte

	if body != nil {
		d, err := json.Marshal(body)
		data = d

		if err != nil {
			return err
		}
	}

	w.WriteHeader(status)
	w.Write(data)
	return nil
}
