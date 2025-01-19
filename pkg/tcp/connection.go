package tcp

import (
	"errors"
	"io"
	"log/slog"
	"math"
	"net"
	"os"
	"sync"
	"time"

	"github.com/hellojonas/flog/pkg/applog"
)

type TCPConnection struct {
	app string

	mu   sync.Mutex
	conn net.Conn
}

func NewTCPConnection(conn net.Conn) *TCPConnection {
	return &TCPConnection{
		conn: conn,
	}
}

func (c *TCPConnection) Conn() net.Conn {
	return c.conn
}

func (c *TCPConnection) App() string {
	return c.app
}

func (c *TCPConnection) RecvTimeout(timout time.Time) ([]byte, error) {
	var data []byte
	logger := applog.Logger().With(slog.String("connection", c.conn.RemoteAddr().String()))
	hChunk := make([]byte, MESSAGE_HEADER_SIZE)

	for {
		c.conn.SetReadDeadline(timout)
		c.mu.Lock()
		n, err := c.conn.Read(hChunk)
		c.mu.Unlock()

		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, os.ErrDeadlineExceeded) {
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
		err = msg.UnmarshalHeaderBinary(hChunk)
		if err != nil {
			return nil, err
		}
		payload := make([]byte, msg.Length)

		c.conn.SetReadDeadline(timout)
		c.mu.Lock()
		n, err = c.conn.Read(payload)
		c.mu.Unlock()

		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, os.ErrDeadlineExceeded) {
				return nil, err
			}
			logger.Error("error reading data from connection", slog.Any("err", err))
			continue
		}

		if n == 0 {
			logger.Warn("read 0 bytes from connection")
			continue
		}

		data = append(data, payload...)

		if msg.Flags&FLAG_MESSAGE_END != 0 {
			return data, nil
		}
	}
}

func (c *TCPConnection) SendWithFlags(data []uint8, flags TCPMessageFlag) error {
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
			end = start + len(data[start:]) // TODO: review
		}

		msg := TCPMessage{
			Flags: msgFlags,
			Data:  data[start:end],
		}

		payLoad, err := msg.MarshalBinary()

		if err != nil {
			return err
		}

		c.conn.SetWriteDeadline(time.Time{})
		c.mu.Lock()
		n, err := c.conn.Write(payLoad)
		c.mu.Unlock()

		if err != nil {
			return errors.New("Send: " + err.Error())
		}

		if n == 0 {
			return errors.New("Send: no data written.")
		}
	}

	return nil
}

func (c *TCPConnection) Recv() ([]byte, error) {
	return c.RecvTimeout(time.Time{})
}

func (c *TCPConnection) Send(data []uint8) error {
	return c.SendWithFlags(data, TCPMessageFlag(0))
}
