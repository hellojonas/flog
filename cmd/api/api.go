package main

import (
	"database/sql"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hellojonas/flog/pkg/applog"
	"github.com/hellojonas/flog/pkg/services"
	"github.com/hellojonas/flog/pkg/services/api"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	addr := ":8080"
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

	mux := http.NewServeMux()

	logsDir := filepath.Join(flogHome, "logs")
	logService := services.NewLogService(db, logsDir)

	appService := services.NewAppService(db)
	appRouter := api.NewAppRouter(appService, logService)
	appRouter.Route(mux)

	usrRouter := api.NewUserRouter(services.NewUserService(db), appService)
	usrRouter.Route(mux)

	err = http.ListenAndServe(addr, mux)

	if err != nil {
		panic(err)
	}
}
