package tcp

import (
	"encoding/binary"
	"errors"
)

const (
	FLAG_PART_CONTINUE = byte(0)
	FLAG_PART_END      = byte(1)

	MESSAGE_VERSION = byte(1)
	HEADER_LENGTH   = 4
)

// Message [ [Version:1] [Flag:1] [Length:2] [Data] ]

type Message struct {
	// Command tells how this message should be parsed
	Flag byte
	Data []byte
}

func (m *Message) MarshalBinary() ([]byte, error) {
	msgLen := len(m.Data) + HEADER_LENGTH

	if msgLen <= HEADER_LENGTH {
		return nil, errors.New("Message#MarshalBinary error: empty message.")
	}

	msg := make([]byte, HEADER_LENGTH)
	msg[0] = MESSAGE_VERSION
	msg[1] = m.Flag
	binary.BigEndian.PutUint16(msg[2:], uint16(len(m.Data)))

	msg = append(msg, m.Data...)

	if len(msg) != msgLen {
		return nil, errors.New("Message#MarshalBinary Error: expected message length differs from actual length.")
	}

	return msg, nil
}

func (m *Message) UnmarshalBinary(msg []byte) error {
	if len(msg) < HEADER_LENGTH {
		return errors.New("Message#UnmarshalBinary error: message length is less than required header length")
	}

	version := msg[0]

	if version != MESSAGE_VERSION {
		return errors.New("Message#UnmarshalBinary error: unsupported message version")
	}

	flag := msg[1]
	length := binary.BigEndian.Uint16(msg[2:])

	data := make([]byte, length)
	copy(data, msg[HEADER_LENGTH:])

	m.Flag = flag
	m.Data = data

	return nil
}
