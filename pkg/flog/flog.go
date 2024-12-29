package flog

import (
	"log/slog"

	"github.com/hellojonas/flog/pkg/applog"
	"github.com/hellojonas/flog/pkg/tcp"
)

type flog struct {
}

func New() *flog {
	return &flog{}
}

func (f *flog) Handle(client *tcp.TCPClient) {
	logger := applog.Logger().With(slog.String("app", client.App()))

	for {
		data, err := client.Recv()

		if err != nil {
			logger.Error("Error reading from client", slog.Any("err", err))
			break
		}

		logger.Info("message received", slog.Int64("len", int64(len(data))))
	}
}
