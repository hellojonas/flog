package applog

import (
	"log/slog"
	"sync"
	"testing"
)

func _TestAppLogger(t *testing.T) {
	loggers := 5
	records := 10000

	wg := sync.WaitGroup{}

	for i := range loggers {
		wg.Add(1)
		go func(service int) {
			logger := Logger()
			for range records {
				logger.Info("Logging an importante message", slog.Int("logger", service))
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}
