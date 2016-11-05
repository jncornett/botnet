// Masters wait for slaves to connect, and then begin issuing commands to them
package main

import (
	"flag"
	"log"
	"net"
	"net/rpc"
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
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		log.Print("accepted connection from", conn.RemoteAddr())
		go func() {
			defer conn.Close()
			client := rpc.NewClient(conn)
			var reply bool
			{
				// Issue commands here
				log.Print("calling Service.SayHello() on", conn.RemoteAddr())
				err = client.Call("Service.SayHello", false, &reply)
				if err != nil {
					log.Print(err)
				}
			}
		}()
	}
}
