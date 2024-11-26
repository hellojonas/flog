package applogger

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

const LOG_FILE_FORMAT = "2006-02-01"

type AppLogHandler struct {
	dest     string
	filename string
	handler  slog.Handler
}

func filename(time.Time) string {
	return time.Now().Format(LOG_FILE_FORMAT) + ".log"
}

var defaultLogger *slog.Logger

func Default() *slog.Logger {
	if defaultLogger != nil {
		return defaultLogger
	}

	h := newAppLogHandler("./logs")
	defalutLogger := slog.New(h)

	return defalutLogger
}

func newAppLogHandler(dest string) *AppLogHandler {

	if err := os.MkdirAll(dest, fs.ModePerm); err != nil {
		panic(err)
	}

	filename := filename(time.Now())

	path := filepath.Join(dest, filename)
	out, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)

	if err != nil {
		panic(err)
	}

	th := slog.NewTextHandler(out, &slog.HandlerOptions{})

	h := &AppLogHandler{
		dest:     dest,
		filename: filename,
		handler:  th,
	}

	return h
}

func (h *AppLogHandler) Handle(ctx context.Context, r slog.Record) error {
	filename := filename(time.Now())

	if h.filename == filename {
		return h.handler.Handle(ctx, r)
	}

	path := filepath.Join(h.dest, filename)
	out, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)

	if err != nil {
		return err
	}

	newHandler := slog.NewTextHandler(out, &slog.HandlerOptions{})

	h.filename = filename
	h.handler = newHandler

	return nil
}

// Handler returns the Handler wrapped by h.
func (h *AppLogHandler) Handler() slog.Handler {
	return h.handler
}

func (h *AppLogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= slog.LevelDebug
}

func (h *AppLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return createLogHandler(h.dest, h.handler.WithAttrs(attrs))
}

func (h *AppLogHandler) WithGroup(name string) slog.Handler {
	return createLogHandler(h.dest, h.handler.WithGroup(name))
}

func createLogHandler(dest string, h slog.Handler) *AppLogHandler {
	if err := os.MkdirAll(dest, fs.ModePerm); err != nil {
		panic(err)
	}

	filename := filename(time.Now())

	path := filepath.Join(dest, filename)
	out, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)

	if err != nil {
		panic(err)
	}

	if lh, ok := h.(*AppLogHandler); ok {
		h = lh.Handler()
	} else {
		h = slog.NewTextHandler(out, &slog.HandlerOptions{})
	}

	return &AppLogHandler{
		dest:     dest,
		filename: filename,
		handler:  h,
	}
}
