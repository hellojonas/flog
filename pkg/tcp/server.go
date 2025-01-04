package tcp

import (
	"errors"
	"log/slog"
	"net"

	"github.com/hellojonas/flog/pkg/applog"
)

type TCPServer struct {
	listener net.Listener
	handler   TCPHandler
}

type TCPHandler interface {
	Handle(*TCPConnection)
}

func NewTCPServer(addr string, handler TCPHandler) (*TCPServer, error) {
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, err
	}

	if handler == nil {
		return nil, errors.New("handler cannot be nil")
	}

	s := &TCPServer{
		listener: listener,
		handler:   handler,
	}

	return s, nil
}

func (ts *TCPServer) Close() error {
	return ts.listener.Close()
}

func (ts *TCPServer) StartAccept() {
	logger := applog.Logger()
	logger.Info("accepting connections at", slog.String("addr", ts.listener.Addr().String()))
	for {
		conn, err := ts.listener.Accept()
		logger.Info("connection established", slog.String("addr", conn.RemoteAddr().String()))
		if err != nil {
			logger.Error("error accepting connection", slog.Any("err", err), slog.String("addr", conn.RemoteAddr().String()))
			continue
		}
		client := NewTCPConnection(conn)
		go ts.handler.Handle(client)
	}
}
