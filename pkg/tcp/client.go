package tcp

import (
	"errors"
	"io"
	"log/slog"
	"math"
	"net"
	"time"

	"github.com/hellojonas/flog/pkg/applog"
)

type TCPClient struct {
	app  string
	conn net.Conn
}

func NewTCPClient(conn net.Conn) *TCPClient {
	return &TCPClient{
		conn: conn,
	}
}

func (c *TCPClient) App() string {
	return c.app
}

func (c *TCPClient) Recv() ([]byte, error) {
	var data []byte
	logger := applog.Logger().With(slog.String("client", c.conn.RemoteAddr().String()))
	chunk := make([]uint8, MESSAGE_MAX_LENGTH)

	for {
		c.conn.SetReadDeadline(time.Time{})
		n, err := c.conn.Read(chunk)

		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, err
			}
			logger.Error("error reading data from connection", slog.Any("err", err))
			continue
		}

		if n == 0 {
			logger.Warn("read 0 bytes from connection")
			continue
		}

		msg := TCPMessage{}
		msg.UnmarshalBinary(chunk)

		data = append(data, msg.Data...)

		if msg.Flags&FLAG_MESSAGE_END != 0 {
			return data, nil
		}
	}
}

func (c *TCPClient) Send(data []uint8) error {
	return c.SendWithFlags(data, TCPMessageFlag(0))
}

func (c *TCPClient) SendWithFlags(data []uint8, flags TCPMessageFlag) error {
	maxPayload := MESSAGE_MAX_LENGTH - MESSAGE_HEADER_SIZE
	parts := int(math.Ceil(float64(len(data)) / float64(maxPayload)))

	for i := range parts {
		start := i * maxPayload
		end := (i + 1) * maxPayload
		msgFlags := flags

		if i == 0 {
			msgFlags |= FLAG_MESSAGE_START
		}

		if i == parts-1 {
			msgFlags |= FLAG_MESSAGE_END
			end = start + len(data[start:])
		}

		msg := TCPMessage{
			Flags: msgFlags,
			Data:  data[start:end],
		}

		payLoad, err := msg.MarshalBinary()

		if err != nil {
			return err
		}

		n, err := c.conn.Write(payLoad)

		if err != nil {
			return errors.New("Send: " + err.Error())
		}

		if n == 0 {
			return errors.New("Send: no data written.")
		}
	}

	return nil
}
