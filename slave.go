package botnet

import (
	"io"
	"log"
	"net/rpc"
	"time"
)

const (
	DefaultSlaveConnectInterval = 5 * time.Second
)

func DefaultSlaveNewServer() Server {
	return rpc.NewServer()
}

type Server interface {
	Register(interface{}) error
	ServeConn(io.ReadWriteCloser)
}

type Slave struct {
	ConnectInterval time.Duration
	server          Server
}

func NewSlave(newServer func() Server) *Slave {
	if newServer == nil {
		newServer = DefaultSlaveNewServer
	}
	return &Slave{DefaultSlaveConnectInterval, newServer()}
}

func (s *Slave) Serve(connect func() (io.ReadWriteCloser, error)) {
	for {
		s.handleConnection(connectLoop(s.ConnectInterval, connect))
	}
}

func (s *Slave) Register(svc interface{}) {
	s.server.Register(svc)
}

func (s *Slave) handleConnection(conn io.ReadWriteCloser) {
	defer conn.Close()
	s.server.ServeConn(conn)
}

func connectLoop(
	interval time.Duration,
	connect func() (io.ReadWriteCloser, error),
) io.ReadWriteCloser {
	for {
		conn, err := connect()
		if err == nil {
			return conn
		}
		log.Print(err)
		time.Sleep(interval)
	}
}
