package logger

import (
	"encoding/json"
	"io"
	"net/http"
)

type controller struct {
	store chan Entry
}

func NewController(entries chan Entry) controller {
	return controller{
		store: entries,
	}
}

func (c controller) LogHandler(w http.ResponseWriter, r *http.Request) {
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
	c.store <- e
	w.WriteHeader(http.StatusOK)
}
