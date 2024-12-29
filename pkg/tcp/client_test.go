package tcp

import (
	"bytes"
	"net"
	"testing"
	"time"
)

func TestClientSendReadMessage(t *testing.T) {
	server, client := net.Pipe()
	defer func() {
		server.Close()
		client.Close()
	}()

	data := make(chan []uint8)
	defer close(data)

	go func(dataC chan []uint8) {
		_data := make([]uint8, 0)
		server.SetDeadline(time.Now().Add(5 * time.Second))
		for {
			payload := make([]uint8, MESSAGE_MAX_LENGTH)
			_, err := server.Read(payload)
			// fmt.Printf("Server Read: %d\n", n)
			if err != nil {
				break
			}
			server.SetDeadline(time.Now().Add(5 * time.Second))
			msg := TCPMessage{}
			msg.UnmarshalBinary(payload)
			_data = append(_data, msg.Data...)
		}
		dataC <- _data
	}(data)

	tcpClient := NewTCPClient(client)
	payloadLen := (12 * MESSAGE_MAX_LENGTH)
	payload := bytes.Repeat([]uint8("A"), payloadLen)

	err := tcpClient.Send(payload)
	tcpClient.conn.Close()

	if err != nil {
		t.Fatal(err)
	}

	select {
	case d := <-data:
		actualLen := len(d)
		if actualLen != payloadLen {
			t.Fatalf("expected payload size: %d, acutal: %d\n", payloadLen, actualLen)
		}
	}
}
