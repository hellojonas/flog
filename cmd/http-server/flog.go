package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hellojonas/flog/pkg/logger"
)

func main() {
	bufSize := flag.Int("bufsize", 100, "Log entries buffers size")
	port := flag.Int("port", 8080, "Port to listen")
	logDir := flag.String("log-dir", "", "Root log directory")
	poll := flag.Int("poll", 700, "Poll interval in ms")
	flag.Parse()
	addr := ":" + strconv.FormatInt(int64(*port), 10)
	store := make(chan logger.Entry, *bufSize)

	fmt.Println(*bufSize, *port, *logDir, addr)

	if logDir == nil || *logDir == "" {
		home := os.Getenv("HOME")
		p := (filepath.Join(home, ".flog", "logs"))
		logDir = &p
	}

	lggr, err := logger.New(*logDir, store)

	if err != nil {
		panic(err)
	}

	go lggr.Log(*poll)
	server := http.NewServeMux()

	server.Handle("/logs", lggr)
	log.Println("application listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, server))
}
