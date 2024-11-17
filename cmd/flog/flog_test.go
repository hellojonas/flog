package main

import (
	"bytes"
	"testing"

	"github.com/hellojonas/flog/pkg/tcp"
)

func TestConnection(t *testing.T) {
	client, err := tcp.NewTCPClient("127.0.0.1", 4004)

	if err != nil {
		t.Fatal(err)
	}

	data := bytes.Repeat([]byte("AB"), 10*512)

	client.Send(data)

	select {}
}
