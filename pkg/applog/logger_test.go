package applog

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"
)

func TestAppLogger(t *testing.T) {
	loggers := 5
	records := 10000

	wg := sync.WaitGroup{}
	start := time.Now()

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

	end := time.Now()
	elapsed := end.UnixMilli() - start.UnixMilli()

}

func TestSlog(t *testing.T) {
	loggers := 5
	records := 10000

	wg := sync.WaitGroup{}
	start := time.Now()

	for i := range loggers {
		wg.Add(1)
		go func(service int) {
			out, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, fs.ModePerm)
			if err != nil {
				t.Fatal(err)
			}
			logger := slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{}))
			for range records {
				logger.Info("Logging an importante message", slog.Int("logger", service))
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	end := time.Now()
	elapsed := end.UnixMilli() - start.UnixMilli()

	fmt.Printf("elapsed: %dms (%.2f)s", elapsed, float64(elapsed)/float64(time.Millisecond))
}
