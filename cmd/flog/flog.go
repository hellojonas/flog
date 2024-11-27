package main

import (
	"os"

	"github.com/hellojonas/flog/pkg/applog"
)

func main() {
	logsDir := os.Getenv("HOME") + "/.local/flog/logs"
	applog.Config(&applog.AppLogOpts{
		Dest: logsDir,
	})
}
