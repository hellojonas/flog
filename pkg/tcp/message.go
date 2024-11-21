package tcp

import (
	"encoding/binary"
	"errors"
)

type TCPMessageFlag uint8

const (
	FLAG_PART_START TCPMessageFlag = 1 << iota
	FLAG_PART_END
	FLAG_AUTH
	FLAG_ERROR
)

const (
	MESSAGE_VERSION = uint8(1)
	HEADER_LENGTH   = 4
)

// Message [ [Version:1] [Flags:1] [Length:2] [Data] ]
// Flags {part_start, part_continue, part_end, error, auth}

type TCPMessage struct {
	// Command tells how this message should be parsed
	Flags TCPMessageFlag // TODO: change this to support multiple flags
	Data  []byte
}

func (m *TCPMessage) MarshalBinary() ([]byte, error) {
	msgLen := len(m.Data) + HEADER_LENGTH

	if msgLen <= HEADER_LENGTH {
		return nil, errors.New("Message#MarshalBinary error: empty message.")
	}

	msg := make([]byte, HEADER_LENGTH)
	msg[0] = MESSAGE_VERSION
	msg[1] = uint8(m.Flags)
	binary.BigEndian.PutUint16(msg[2:], uint16(len(m.Data)))

	msg = append(msg, m.Data...)

	if len(msg) != msgLen {
		return nil, errors.New("Message#MarshalBinary Error: expected message length differs from actual length.")
	}

	return msg, nil
}

func (m *TCPMessage) UnmarshalBinary(msg []byte) error {
	if len(msg) < HEADER_LENGTH {
		return errors.New("Message#UnmarshalBinary error: message length is less than required header length")
	}

	version := msg[0]

	if version != MESSAGE_VERSION {
		return errors.New("Message#UnmarshalBinary error: unsupported message version")
	}

	flags := msg[1]
	length := binary.BigEndian.Uint16(msg[2:])

	data := make([]byte, length)
	copy(data, msg[HEADER_LENGTH:])

	m.Flags = TCPMessageFlag(flags)
	m.Data = data

	return nil
}
