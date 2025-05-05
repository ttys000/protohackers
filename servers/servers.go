package servers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type LocalTCPServer struct {
	Name                     string
	Port                     int
	ConnectionHandler        func(net.Conn)
	MaxConcurrentConnections int
	ConnectionTracker        *sync.WaitGroup
}

func (server *LocalTCPServer) Serve(ctx context.Context) {
	tcp, err := net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", server.Port))
	if err != nil {
		log.Fatalf("%s Can not listen on %d/tcp: %s", server.Name, server.Port, err)
	}
	log.Printf("%s Server started, listening on %s \n", server.Name, tcp.Addr())

	go func() {
		<-ctx.Done()
		log.Printf("%s Shutdown signal received, stopping new connections...", server.Name)
		tcp.Close()
	}()

	server.handleTCP(ctx, tcp)
}

func (server *LocalTCPServer) handleTCP(ctx context.Context, tcp net.Listener) {
	semaphore := make(chan struct{}, server.MaxConcurrentConnections)

	for {
		select {
		case <-ctx.Done():
			log.Printf("%s Stopping new connection acceptance", server.Name)
			return
		case semaphore <- struct{}{}:
			// Try to accept a new connection
			conn, err := tcp.Accept()
			if err != nil {
				<-semaphore // Release the semaphore since we failed to accept
				if errors.Is(err, net.ErrClosed) {
					log.Printf("%s Listener closed, stopping connection acceptance", server.Name)
					return
				}
				log.Printf("%s Accepting connection error: %v", server.Name, err)
				continue
			}

			// Set a read deadline to prevent hanging connections
			if err := conn.SetReadDeadline(time.Now().Add(5 * time.Minute)); err != nil {
				log.Printf("%s Failed to set read deadline: %v", server.Name, err)
				conn.Close()
				<-semaphore
				continue
			}

			server.ConnectionTracker.Add(1)
			go func() {
				defer func() {
					conn.Close()
					<-semaphore
					server.ConnectionTracker.Done()
				}()
				server.ConnectionHandler(conn)
			}()
		default:
			// We can't acquire a semaphore slot, so log the rejection
			log.Printf("%s Connection rejected: maximum concurrent connections reached", server.Name)
		}
	}
}
