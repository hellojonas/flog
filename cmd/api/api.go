package main

import (
	"database/sql"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hellojonas/flog/pkg/applog"
	"github.com/hellojonas/flog/pkg/apps"
	"github.com/hellojonas/flog/pkg/users"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	addr := ":8080"
	logger := applog.Logger()
	userDir := os.Getenv("HOME")

	if os.PathSeparator == '\\' {
		userDir = os.Getenv("USERPROFILE")
	}

	flogRoot := filepath.Join(userDir, ".flog")

	if err := os.MkdirAll(flogRoot, fs.ModePerm); err != nil {
		logger.Error("error creating flog root dir.", slog.Any("err", err))
		panic(err)
	}

	dbpath := "file:" + filepath.Join(flogRoot, "flog.db")

	logger.Info("opening db...", slog.String("db", dbpath))
	db, err := sql.Open("sqlite3", dbpath)

	if err != nil {
		logger.Error("error opening database connection.", slog.Any("err", err))
		panic(err)
	}

	defer db.Close()

	mux := http.NewServeMux()

	appService := apps.NewService(db)
	appRouter := apps.NewRouter(appService)
	appRouter.Route(mux)

	usrRouter := users.NewRouter(users.NewService(db), appService)
	usrRouter.Route(mux)

	err = http.ListenAndServe(addr, mux)

	if err != nil {
		panic(err)
	}
}
