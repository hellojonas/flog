package tcp

import (
	"encoding/binary"
	"testing"
)

func TestMarshalBinary(t *testing.T) {
	m := TCPMessage{
		Flags: FLAG_PART_END,
		Data: []byte("hello!"),
	}

	msg, err := m.MarshalBinary()

	if err != nil {
		t.Fatal(err)
	}

	msgVer := msg[0]

	if msgVer != MESSAGE_VERSION {
		t.Fatalf("expected message version: %d, actual: %d\n", MESSAGE_LENGTH, msgVer)
	}

	flag := msg[1]
	if flag != FLAG_PART_END {
		t.Fatalf("expected message flag: %d, actual: %d\n", FLAG_PART_END, flag)
	}

	dataLen := binary.BigEndian.Uint16(msg[2:])
	if int(dataLen) != len(m.Data) {
		t.Fatalf("expected message data length: %d, actual: %d\n", len(m.Data), dataLen)
	}
}

func TestUnmarshalBinary(t *testing.T) {
	m := TCPMessage{
		Flags: FLAG_PART_END,
		Data: []byte("hello!"),
	}

	msg, err := m.MarshalBinary()

	if err != nil {
		t.Fatal(err)
	}

	m2 := TCPMessage{}

	err = m2.UnmarshalBinary(msg)

	if err != nil {
		t.Fatal(err)
	}

	if m2.Flags != m.Flags {
		t.Fatalf("expected flag: %d, actual: %d\n", m.Flags, m2.Flags)
	}

	if len(m2.Data) != len(m.Data) {
		t.Fatalf("expected message data length: %d, actual: %d\n", len(m.Data), len(m2.Data))
	}

	for i := HEADER_LENGTH; i < len(m2.Data); i++ {
		if m2.Data[i] == m.Data[i] {
			continue
		}

		t.Fatalf("expected %d at data index %d, actual: %d\n", m.Data[i], i, m2.Data[i])
	}
}
