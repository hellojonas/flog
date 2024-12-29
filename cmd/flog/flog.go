package main

import (
	"os"
	"path"

	"github.com/hellojonas/flog/pkg/applog"
	"github.com/hellojonas/flog/pkg/flog"
	"github.com/hellojonas/flog/pkg/tcp"
)

func main() {
	logsDir := ""
	userDir := os.Getenv("HOME")
	addr := ":8008"

	if os.PathSeparator == '\\' {
		userDir = os.Getenv("USERPROFILE")
	}

	logsDir = path.Join(userDir, ".flog", "logs")

	applog.Config(&applog.AppLogOpts{
		Dest: logsDir,
	})

	server, err := tcp.NewTCPServer(addr, flog.New())

	if err != nil {
		panic(err)
	}

	server.StartAccept()
}
