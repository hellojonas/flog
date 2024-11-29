package tcp

import "testing"

func TestMarshalBinary(t *testing.T) {
	text := "Hello there!"
	msg := TCPMessage{
		Flags: FLAG_MESSAGE_START | FLAG_MESSAGE_END,
		Data:  []uint8(text),
	}

	if _, err := msg.MarshalBinary(); err != nil {
		t.Fatal(err)
	}
}

func TestUnmarshalBinary(t *testing.T) {
	text := "Hello there!"
	msg := TCPMessage{
		Flags: FLAG_MESSAGE_START | FLAG_MESSAGE_END,
		Data:  []uint8(text),
	}

	data, err := msg.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	actual := TCPMessage{}
	if err := actual.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}

	actualText := string(actual.Data)

	if actual.Flags != msg.Flags {
		t.Fatalf("expected flags: '%08b', actual flags: '%08b'\n", msg.Flags, actual.Flags)
	}

	if actualText != text {
		t.Fatalf("expected message: '%s', actual message: '%s'\n", text, actualText)
	}
}
