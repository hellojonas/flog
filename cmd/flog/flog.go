package main

import (
	"github.com/hellojonas/flog/pkg/flog"
	"github.com/hellojonas/flog/pkg/tcp"
)

func main() {
	addr := ":8008"
	server, err := tcp.NewTCPServer(addr, flog.New())

	if err != nil {
		panic(err)
	}

	defer server.Close()

	server.StartAccept()
}
