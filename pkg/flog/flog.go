package flog

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/hellojonas/flog/pkg/applog"
	"github.com/hellojonas/flog/pkg/tcp"
)

type flog struct {
	appId string
}

type authInfo struct {
	AppId  string `json:"app_id"`
	Secret string `json:"secret"`
}

func New() *flog {
	return &flog{}
}

func (f *flog) authenticate(client *tcp.TCPConnection) error {
	err := client.SendWithFlags([]byte("AUTH_REQUEST"), tcp.FLAG_MESSAGE_AUTH)
	conn := client.Conn()

	if err != nil {
		return err
	}

	chunk := make([]byte, tcp.MESSAGE_MAX_LENGTH)
	conn.SetReadDeadline(time.Now().Add(15 * time.Second))
	n, err := conn.Read(chunk)

	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("0 bytes read from connection. auth failed")
	}

	authMsg := tcp.TCPMessage{}
	err = authMsg.UnmarshalBinary(chunk)

	if err != nil {
		return err
	}

	if authMsg.Flags&tcp.FLAG_MESSAGE_AUTH == 0 {
		return errors.New("not authenticated.")
	}

	if authMsg.Flags&tcp.FLAG_MESSAGE_START == 0 &&
		authMsg.Flags&tcp.FLAG_MESSAGE_END == 0 {
		fmt.Printf("%08b\n", authMsg.Flags)
		return errors.New("auth info must be contained in a single message")
	}

	var ai authInfo
	err = json.Unmarshal(authMsg.Data, &ai)

	if err != nil {
		return err
	}

	// TODO: validate secret

	err = client.SendWithFlags([]byte("AUTH_OK"), tcp.FLAG_MESSAGE_AUTH)

	if err != nil {
		return err
	}

	f.appId = ai.AppId

	return nil
}

func (f *flog) Handle(client *tcp.TCPConnection) {
	logger := applog.Logger().With(slog.String("app", client.App()))

	// TODO: handle authentication errors logging (general and timeout)
	err := f.authenticate(client)

	if err != nil {
		logger.Error("Error authenticating client", slog.Any("err", err))
		return
	}

	for {
		data, err := client.Recv()

		if err != nil {
			logger.Error("Error reading from client", slog.Any("err", err))
			break
		}

		logger.Info("message received", slog.Int("len", len(data)))
	}
}
