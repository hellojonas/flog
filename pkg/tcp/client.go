package tcp

import (
	"errors"
	"log/slog"
	"math"
	"net"
	"strconv"
)

type TCPClient struct {
	appId   string // TODO: fill this when authenticating
	appName string // TODO: fill this when authenticating
	conn    net.Conn
}

func (c TCPClient) AppId() string {
	return c.appId
}

func (c TCPClient) AppName() string {
	return c.appName
}

func NewTCPClient(host string, port int) (*TCPClient, error) {
	client := TCPClient{}

	addr := host + ":" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	client.conn = conn

	return &client, nil
}

func (c *TCPClient) Send(msg []byte) error {
	return c.sendMessage(msg, TCPMessageFlag(0))
}

func (c *TCPClient) SendError(msg []byte) error {
	return c.sendMessage(msg, FLAG_ERROR)
}

func (c *TCPClient) sendMessage(msg []byte, flags TCPMessageFlag) error {
	if len(msg) < HEADER_LENGTH {
		return errors.New("TCPClient#Send error: message length is less than required header length")
	}

	dataLen := MESSAGE_LENGTH - HEADER_LENGTH
	parts := int(math.Ceil(float64(len(msg)) / float64(dataLen)))

	slog.Info("TCPClient#Send: sending message", "parts", parts, "length", len(msg))

	for i := 0; i < parts; i++ {
		start := dataLen * i
		end := dataLen * (i + 1)

		if i == 0 {
			flags |= FLAG_PART_START
		} else {
			flags &^= FLAG_PART_START
		}

		if i == 0 && i == parts-1 {
			flags = FLAG_PART_START | FLAG_PART_END
		}

		m := TCPMessage{
			Flags: flags,
			Data:  msg[start:end],
		}

		msgBytes, err := m.MarshalBinary()

		if err != nil {
			return err
		}

		n, err := c.conn.Write(msgBytes)

		if err != nil {
			slog.Info("TCPClient#Send error: could not write message part.", "err", err.Error())
			return errors.New("TCPClient#Send error: could not write message part." + err.Error())
		}

		if n == 0 {
			slog.Info("TCPClient#Send error: no message part written.")
			return errors.New("TCPClient#Send error: no message part written.")
		}
	}

	slog.Info("TCPClient#Send: message sent.", slog.Any("parts", parts))

	return nil
}
