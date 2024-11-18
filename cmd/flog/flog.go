package main

import (
	"github.com/hellojonas/flog/pkg/tcp"
)

func main() {
	server := tcp.NewTCPServer("127.0.0.1", 4004)
	err := server.Listen()

	if err != nil {
		panic(err)
	}

	defer server.Close()

	server.StartAccept()

    if err != nil {
        panic(err)
    }
}
