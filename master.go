package botnet

import (
	"io"
	"log"
	"net/rpc"
)

func DefaultMasterNewClient(conn io.ReadWriteCloser) Client {
	return rpc.NewClient(conn)
}

type Client interface {
	Call(string, interface{}, interface{}) error
}

type Master struct {
	NewClient func(io.ReadWriteCloser) Client
}

func NewMaster() *Master {
	return &Master{DefaultMasterNewClient}
}

func (m *Master) Serve(
	listen func() (io.ReadWriteCloser, error),
	command func(Client),
) {
	for {
		conn, err := listen()
		if err != nil {
			log.Print(err)
			continue
		}
		go func() {
			defer conn.Close()
			command(m.NewClient(conn))
		}()
	}
}
