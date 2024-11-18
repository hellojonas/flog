package tcp

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"strconv"
	"time"
)

const (
	MESSAGE_LENGTH = 1024
)

type ConnHandler func(net.Conn)

type TCPServer struct {
	host     string
	port     int
	listener net.Listener
	Messages chan TCPMessage
}

func NewTCPServer(host string, port int) TCPServer {
	return TCPServer{
		host: host,
		port: port,
	}
}

func (s *TCPServer) Close() error {
	return s.listener.Close()
}

func (s *TCPServer) Listen() error {
	addr := s.host + ":" + strconv.Itoa(s.port)

	slog.Info("TCP#Listen: starting server...")
	listener, err := net.Listen("tcp", addr)
	slog.Info("TCP#Listen: server started. listening on", slog.Any("addr", addr))

	if err != nil {
		slog.Error("TCP#Listen: could not start server.", slog.Any("addr", addr))
		return errors.New("TCP#Listen: " + err.Error())
	}

	s.listener = listener

	return nil
}

func (s *TCPServer) StartAccept() error {
	slog.Info("TCP#StartAccept: waiting for connection...")
	for {
		conn, err := s.listener.Accept()
		slog.Info("TCP#StartAccept: connection established.", "addr", conn.RemoteAddr())

		if err != nil {
			slog.Error("TCP#StartAccept: cold not establish connection", "err", err, "addr", conn.RemoteAddr())
			continue
		}

		c := TCPClient{
			conn: conn,
		}

		go s.handleConn(&c)
	}
}

func (s *TCPServer) handleConn(client *TCPClient) {
	conn := client.conn
	data := make([]byte, 0)


	for {
		conn.SetDeadline(time.Time{})
		chunk := make([]byte, MESSAGE_LENGTH)
		n, err := conn.Read(chunk)

		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Error("StartAccept#handleConn: EOF.")
				return
			}

			slog.Error("StartAccept#handleConn: could not read data from connection.", "err", err)
			return
		}

		if n == 0 {
			slog.Error("StartAccept#handleConn: nothing read.")
			continue
		}

		msg := TCPMessage{}
		err = msg.UnmarshalBinary(chunk)

		if err != nil {
			// TODO: what if ono of the parts fails to be delivered? handle here
			// TODO: write an error to the client to signal the error and let im handle?
			// TODO: Add more control fields on message header or use one control field with bitwise operators (start, continue, end, error)
			slog.Error("StartAccept#handleConn: could not parse message.", "err", err)
		}

		data = append(data, msg.Data...)

		if msg.Flag == FLAG_PART_END {
			s.Messages <- msg
			slog.Info("StartAccept#handleConn: message received.", "length", len(data))
			continue
		}
	}
}
