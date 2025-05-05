package echo

import (
	"context"
	"net"
	"sync"

	"protohackers/servers"
)

const serverName = "[Echo]"

func RunServer(ctx context.Context, port int, wg *sync.WaitGroup) {
	server := &servers.LocalTCPServer{
		Name:                     serverName,
		Port:                     port,
		ConnectionHandler:        handleConn,
		MaxConcurrentConnections: 100,
		ConnectionTracker:        wg,
	}
	server.Serve(ctx)
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
