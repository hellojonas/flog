package tcp

import (
	"encoding/binary"
	"errors"
)

type TCPMessage struct {
	Flags  TCPMessageFlag
	Length uint16
	Data   []uint8
}

type TCPMessageFlag uint8

const (
	FLAG_MESSAGE_START TCPMessageFlag = 1 << iota
	FLAG_MESSAGE_END
	FLAG_MESSAGE_ERROR
	FLAG_MESSAGE_AUTH
)

const (
	MESSAGE_VERSION     uint8 = 1
	MESSAGE_HEADER_SIZE       = 4
	MESSAGE_MAX_LENGTH        = 1024
)

// Message format
// [[VERSION:1][FLAGS:1][LENGTH:2][DATA:n]]

func (m *TCPMessage) MarshalBinary() ([]uint8, error) {
	dataLen := len(m.Data)

	if dataLen == 0 {
		return nil, errors.New("MarshalBinary: empty message")
	}

	data := make([]uint8, MESSAGE_HEADER_SIZE)
	data[0] = MESSAGE_VERSION
	data[1] = uint8(m.Flags)
	binary.BigEndian.PutUint16(data[2:], uint16(dataLen))
	data = append(data, m.Data...)

	return data, nil
}

func (m *TCPMessage) UnmarshalHeaderBinary(data []uint8) error {
	version := data[0]

	if version != MESSAGE_VERSION {
		return errors.New("UnmarshalBinary: unsupported message version")
	}

	if len(data) < MESSAGE_HEADER_SIZE {
		return errors.New("UnmarshalBinay: invalid message length")
	}

	flags := data[1]
	dataLen := binary.BigEndian.Uint16(data[2:])

	m.Flags = TCPMessageFlag(flags)
	m.Length = dataLen

	return nil
}

func (m *TCPMessage) UnmarshalBinary(data []uint8) error {
	m.UnmarshalHeaderBinary(data)
	d := make([]uint8, m.Length)

	copy(d, data[MESSAGE_HEADER_SIZE:m.Length+MESSAGE_HEADER_SIZE])

	m.Data = d

	return nil
}
