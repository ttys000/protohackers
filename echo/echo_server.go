package echo

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"protohackers/servers"
)

const serverName = "[Echo]"

func RunServer(ctx context.Context, port int, wg *sync.WaitGroup) {
	tcp, err := net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("%s Can not listen on %d/tcp: %s", serverName, port, err)
	}
	log.Printf("%s Server started, listening on %s \n", serverName, tcp.Addr())

	go func() {
		<-ctx.Done()
		log.Printf("%s Shutdown signal received, stopping new connections...", serverName)
		tcp.Close()
	}()

	config := servers.TCPServerConfig{
		ServerName:               serverName,
		ConnectionHandler:        handleConn,
		MaxConcurrentConnections: 100,
	}
	servers.HandleTCP(ctx, tcp, wg, config)
}

func handleConn(conn net.Conn) {
	buf := make([]byte, 65536)
	for {
		nbyte, err := conn.Read(buf)
		if err != nil || nbyte == 0 {
			return
		}
		conn.Write(buf[:nbyte])
	}
}
