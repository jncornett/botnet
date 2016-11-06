// Slaves call back to the master and then accept arbitrary payloads
package main

import (
	"flag"
	"io"
	"log"
	"net"

	"github.com/jncornett/botnet"
)

const (
	defaultMaster = "localhost:13000"
)

func main() {
	var (
		master = flag.String("master", defaultMaster, "master node")
	)
	flag.Parse()
	service := Service{}
	slave := botnet.NewSlave(nil)
	slave.Register(&service)
	slave.Serve(func() (io.ReadWriteCloser, error) { return net.Dial("tcp", *master) })
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
