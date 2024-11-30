package tcp

import (
	"errors"
	"math"
	"net"
)

type TCPClient struct {
	conn net.Conn
}

func NewWithConnection(conn net.Conn) *TCPClient {
	return &TCPClient{
		conn: conn,
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
