package flog

import (
	"encoding/json"
	"errors"
	"net"
	"time"

	"github.com/hellojonas/flog/pkg/tcp"
)

type client struct {
	app  string
	conn tcp.TCPConnection
	data chan []byte
}

func isAuthMsg(msg *tcp.TCPMessage) bool {
	return msg.Flags&tcp.FLAG_MESSAGE_AUTH != 0 &&
		msg.Flags&tcp.FLAG_MESSAGE_START != 0 &&
		msg.Flags&tcp.FLAG_MESSAGE_END != 0
}

func (c *client) authenticate(cred ClientCredential) error {
	conn := c.conn.Conn()

	chunk := make([]byte, tcp.MESSAGE_MAX_LENGTH)
	conn.SetDeadline(time.Now().Add(AUTH_MESSAGE_TIMEOUT * time.Second))
	n, err := conn.Read(chunk)

	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("0 bytes read from connection. auth failed")
	}

	msg := tcp.TCPMessage{}
	err = msg.UnmarshalBinary(chunk)

	if err != nil {
		return err
	}

	if msg.Flags&tcp.FLAG_MESSAGE_AUTH == 0 {
		return errors.New("not authenticated.")
	}

	if msg.Flags&tcp.FLAG_MESSAGE_START == 0 &&
		msg.Flags&tcp.FLAG_MESSAGE_END == 0 {
		return errors.New("auth message must be contained in a single message")
	}

	if string(msg.Data) != AUTH_REQUEST {
		return errors.New("authentication not requested")
	}

	ccData, err := json.Marshal(cred)

	if err != nil {
		return err
	}

	err = c.conn.SendWithFlags(ccData, tcp.FLAG_MESSAGE_AUTH)

	if err != nil {
		return err
	}

	conn.SetReadDeadline(time.Now().Add(AUTH_MESSAGE_TIMEOUT * time.Second))
	n, err = conn.Read(chunk)

	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("0 bytes read from connection. auth failed")
	}

	err = msg.UnmarshalBinary(chunk)

	if err != nil {
		return err
	}

	if msg.Flags&tcp.FLAG_MESSAGE_AUTH == 0 {
		return errors.New("not authenticated.")
	}

	if msg.Flags&tcp.FLAG_MESSAGE_START == 0 &&
		msg.Flags&tcp.FLAG_MESSAGE_END == 0 {
		return errors.New("auth message must be contained in a single message")
	}

	if string(msg.Data) != AUTH_OK {
		return errors.New("auth failed. auth not ok.")
	}

	return nil
}

func NewClient(addr string, cred ClientCredential) (*client, error) {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	c := client{
		conn: *tcp.NewTCPConnection(conn),
	}

	err = c.authenticate(cred)

	if err != nil {
		return nil, err
	}

	return &c, nil
}
