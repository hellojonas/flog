package logger

import (
	"encoding/json"
	"io"
	"net/http"
)

func (l logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		msg, _ := json.Marshal(map[string]string{"message": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}
	e := Entry{}
	err = json.Unmarshal(body, &e)
	if err != nil {
		msg, _ := json.Marshal(map[string]string{"message": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}
	// TODO: validate log entry
	l.entries <- e
	w.WriteHeader(http.StatusOK)
}
