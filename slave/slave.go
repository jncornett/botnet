// Slaves call back to the master and then accept arbitrary payloads
package main

import (
	"flag"
	"log"
	"net"
	"net/rpc"
	"time"
)

const (
	defaultMaster          = "localhost:13000"
	defaultConnectInterval = 10 * time.Second
)

func main() {
	var (
		master          = flag.String("master", defaultMaster, "master node")
		connectInterval = flag.Duration("interval", defaultConnectInterval, "interval to attempt connecting to master")
	)
	flag.Parse()
	service := Service{}
	server := rpc.NewServer()
	server.Register(&service)
	for {
		func() {
			conn := connect(*connectInterval, *master)
			defer conn.Close()
			server.ServeConn(conn)
		}()
	}
}

func connect(interval time.Duration, addr string) net.Conn {
	for {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Print(err)
		} else {
			return conn
		}
		time.Sleep(interval)
	}
}

type Service struct {
	count int
}

// Add exposed services here

func (s *Service) SayHello(_ bool, _ *bool) error {
	s.count++
	log.Println("Hello", s.count)
	return nil
}
