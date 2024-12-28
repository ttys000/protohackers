package main

import (
	"errors"
	"fmt"
	"log"
	"net"
)

var port = 5002

func main() {
	udp, err := net.ListenPacket("udp", fmt.Sprintf("fly-global-services:%d", port))
	if err != nil {
		log.Fatalf("can't listen on %d/udp: %s", port, err)
	}

	tcp, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("can't listen on %d/tcp: %s", port, err)
	}

	go handleTCP(tcp)

	handleUDP(udp)
}

func handleTCP(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}

			log.Printf("error accepting on %d/tcp: %s", port, err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	buf := make([]byte, 65536)
	for {
		nbyte, err := c.Read(buf)
		if err != nil || nbyte == 0 {
			return
		}

		c.Write(buf)
	}
}

func handleUDP(c net.PacketConn) {
	packet := make([]byte, 2000)

	for {
		n, addr, err := c.ReadFrom(packet)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}

			log.Printf("error reading on %d/udp: %s", port, err)
			continue
		}

		c.WriteTo(packet[:n], addr)
	}
}
