// Masters wait for slaves to connect, and then begin issuing commands to them
package main

import (
	"flag"
	"io"
	"log"
	"net"

	"github.com/jncornett/botnet"
)

const (
	defaultListen = ":13000"
)

func main() {
	var (
		listen = flag.String("listen", defaultListen, "address to listen on")
	)
	flag.Parse()
	log.Print("listening on", *listen)
	ln, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatal(err)
	}
	master := botnet.NewMaster()
	master.Serve(
		func() (io.ReadWriteCloser, error) { return ln.Accept() },
		func(client botnet.Client) {
			var reply bool
			err = client.Call("Service.SayHello", false, &reply)
			if err != nil {
				log.Print(err)
			}
		},
	)
}
