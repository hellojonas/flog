package tcp

import (
	"encoding/json"
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

type TCPServer struct {
	host     string
	port     int
	listener net.Listener
	Messages chan TCPMessage
}

type TCPAuth struct {
	Key string `json:"key"`
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

		c, err := s.authenticate(conn)

		if err != nil {
			slog.Error("TCP#StartAccept: cold not authenticate connection.", "err", err, "addr", conn.RemoteAddr())
			continue
		}

		go s.handleConn(c)
	}
}

func (s *TCPServer) authenticate(conn net.Conn) (*TCPClient, error) {
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	chunk := make([]byte, MESSAGE_LENGTH)
	n, err := conn.Read(chunk)

	if err != nil {
		slog.Error("authenticate: could not read data from connection.", "err", err)
		return nil, err
	}

	if n == 0 {
		slog.Error("StartAccept#authenticate: nothing read.")
		return nil, errors.New("authenticate: nothing read from conection")
	}

	msg := TCPMessage{}
	err = msg.UnmarshalBinary(chunk) // json with application id and application key

	// TODO: authenticate connection here

	if err != nil {
		errors.New("authenticate: could not parse message.")
		slog.Error("StartAccept#authenticate: could not parse message.", "err", err)
	}

	var auth TCPAuth
	err = json.Unmarshal(msg.Data, &auth)

	if err != nil {
		slog.Error("StartAccept#authenticate: could not parse message.", "err", err)
	}

	// reoslve this fields using application parsed auth key
	appId := "demo_app"
	appName := "Demo app"

	c := TCPClient{
		// TODO: add app name here
		conn:    conn,
		appId:   appId,
		appName: appName,
	}

	conn.SetReadDeadline(time.Time{})
	return &c, nil
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

		if msg.Flags & FLAG_PART_END != 0 {
			s.Messages <- msg
			slog.Info("StartAccept#handleConn: message received.", "length", len(data))
			continue
		}
	}
}
