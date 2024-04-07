package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestController(t *testing.T) {
	wg := sync.WaitGroup{}
	errc := make(chan error)
	resc := make(chan *http.Response)
	resCount := 0

	go func() {
		if err := <-errc; err != nil {
			panic(err)
		}
	}()

	go func() {
		for r := range resc {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				errc <- err
			}
            fmt.Println(b)
			wg.Done()
			resCount++
		}
	}()

	reqNum := 500
	for range reqNum {
		wg.Add(1)
		go func(errc chan error, resc chan *http.Response) {
			body, err := json.Marshal(map[string]any{
				"group":   "flog",
				"time":    time.Now(),
				"level":   0,
				"message": "Level 0 log entry",
			})
			if err != nil {
				errc <- err
			}
			r, err := http.NewRequest(http.MethodPost, "http://localhost:8080/logs", bytes.NewReader(body))
			if err != nil {
				errc <- err
			}
			c := http.Client{}
			res, err := c.Do(r)
			if err != nil {
				errc <- err
			}
			resc <- res
		}(errc, resc)
	}
	wg.Wait()

	if reqNum != resCount {
		t.Fatalf("expected %d responses, got %d responses\n", reqNum, resCount)
	}
}
