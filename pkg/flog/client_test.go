package flog

import (
	"fmt"
	"sync"
	"testing"
)

func _TestClientSendMessage(t *testing.T) {
	addr := ":8008"
	cc := ClientCredential{
		AppId:  "flog",
		Secret: "flog_secret_key",
	}

	c, err := NewClient(addr, cc)

	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}

	for i := range 10 {
		wg.Add(1)
		go func() {
			for j := range 100 {
				c.TcpConn.Send([]byte(fmt.Sprintf("Message entry No. %02d - %03d\n", i, j)))
			}
			wg.Done()
		}()
	}

	wg.Wait()
}