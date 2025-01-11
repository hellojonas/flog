package flog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/hellojonas/flog/pkg/applog"
	"github.com/hellojonas/flog/pkg/apps"
	"github.com/hellojonas/flog/pkg/tcp"
	"github.com/hellojonas/flog/pkg/users"
)

const (
	AUTH_OK      = "AUTH_OK"
	AUTH_REQUEST = "AUTH_REQUEST"

	AUTH_MESSAGE_TIMEOUT = 5

	LOG_FILE_FORMAT = "2006-01-02"
)

type flog struct {
	appId   string
	logFile string
	logDir  string
	output  *os.File
	userSvc *users.UserService
	appSvc  *apps.AppService
}

type ClientCredential struct {
	AppId  string `json:"app_id"`
	Secret string `json:"secret"`
}

func New(userSvc *users.UserService, appSvc *apps.AppService) *flog {
	return &flog{
		userSvc: userSvc,
		appSvc:  appSvc,
	}
}

func (f *flog) authenticate(client *tcp.TCPConnection) error {
	err := client.SendWithFlags([]byte("AUTH_REQUEST"), tcp.FLAG_MESSAGE_AUTH)
	conn := client.Conn()

	if err != nil {
		return err
	}

	chunk := make([]byte, tcp.MESSAGE_MAX_LENGTH)
	conn.SetReadDeadline(time.Now().Add(AUTH_MESSAGE_TIMEOUT * time.Second))
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

	var cc ClientCredential
	err = json.Unmarshal(authMsg.Data, &cc)

	if err != nil {
		return err
	}

	app, err := f.appSvc.FindByName(cc.AppId)

	if err != nil {
		return err
	}

	if app.Token != cc.Secret {
		return errors.New("invalid token")
	}

	err = client.SendWithFlags([]byte("AUTH_OK"), tcp.FLAG_MESSAGE_AUTH)


	if err != nil {
		return err
	}

	f.appId = cc.AppId

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

	userDir := os.Getenv("HOME")

	if os.PathSeparator == '\\' {
		userDir = os.Getenv("USERPROFILE")
	}

	dest := filepath.Join(userDir, ".flog", "logs", f.appId)

	f.logDir = dest

	if err := os.MkdirAll(dest, fs.ModePerm); err != nil {
		logger.Error("Error creating logs directory for client", slog.Any("err", err))
		return
	}

	for {
		data, err := client.Recv()

		if err != nil {
			logger.Error("Error reading from client", slog.Any("err", err))
			if errors.Is(err, io.EOF) || errors.Is(err, os.ErrDeadlineExceeded) {
				break
			}
			continue
		}

		err = f.persist(data)

		if err != nil {
			logger.Error("Error persisting data", slog.Any("err", err))
		}

	}
}

func (f *flog) persist(data []byte) error {
	filename := logFilename(time.Now())

	if filename != f.logFile {
		f.logFile = filename
		path := filepath.Join(f.logDir, filename)
		out, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)

		if err != nil {
			return err
		}

		f.output.Close()
		f.output = out
	}

	n, err := f.output.Write(data)

	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("0 bytes writtern to disk")
	}

	return nil
}

func logFilename(time.Time) string {
	return time.Now().Format(LOG_FILE_FORMAT) + ".log"
}
