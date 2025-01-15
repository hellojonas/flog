package main

import (
	"database/sql"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/hellojonas/flog/pkg/applog"
	"github.com/hellojonas/flog/pkg/flog"
	"github.com/hellojonas/flog/pkg/migration"
	"github.com/hellojonas/flog/pkg/services"
	"github.com/hellojonas/flog/pkg/tcp"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	logger := applog.Logger()
	userDir := os.Getenv("HOME")

	if os.PathSeparator == '\\' {
		userDir = os.Getenv("USERPROFILE")
	}

	flogHome := filepath.Join(userDir, ".flog")

	if err := os.MkdirAll(flogHome, fs.ModePerm); err != nil {
		logger.Error("error creating flog root dir.", slog.Any("err", err))
		panic(err)
	}

	dbpath := "file:" + filepath.Join(flogHome, "flog.db")

	logger.Info("opening db...", slog.String("db", dbpath))
	db, err := sql.Open("sqlite3", dbpath)

	if err != nil {
		logger.Error("error opening database connection.", slog.Any("err", err))
		panic(err)
	}

	defer db.Close()

	migrationPath := filepath.Join("migrations")
	err = migration.Migrate(db, migrationPath)

	if err != nil {
		logger.Error("error migrating flog database.", slog.Any("err", err))
		panic(err)
	}

	addr := ":8008"
	logDir := filepath.Join(flogHome, "logs")
	appSvc := services.NewAppService(db)
	logSvc := services.NewLogService(db, logDir)

	flogHandler := flog.New(appSvc, logSvc, flogHome)
	server, err := tcp.NewTCPServer(addr, flogHandler)

	if err != nil {
		logger.Error("error staring server.", slog.Any("err", err))
		panic(err)
	}

	defer server.Close()

	server.StartAccept()
}
