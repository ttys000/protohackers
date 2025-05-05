package servers

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"
)

const maxConcurrentConnections = 100

type TCPServerConfig struct {
	ServerName        string
	ConnectionHandler func(net.Conn)
}

func HandleTCP(ctx context.Context, tcp net.Listener, wg *sync.WaitGroup, config TCPServerConfig) {
	semaphore := make(chan struct{}, maxConcurrentConnections)

	for {
		select {
		case <-ctx.Done():
			log.Printf("%s Stopping new connection acceptance", config.ServerName)
			return
		case semaphore <- struct{}{}:
			// Try to accept a new connection
			conn, err := tcp.Accept()
			if err != nil {
				<-semaphore // Release the semaphore since we failed to accept
				if errors.Is(err, net.ErrClosed) {
					log.Printf("%s Listener closed, stopping connection acceptance", config.ServerName)
					return
				}
				log.Printf("%s Accepting connection error: %v", config.ServerName, err)
				continue
			}

			wg.Add(1)
			go func() {
				defer func() {
					<-semaphore
					wg.Done()
				}()
				config.ConnectionHandler(conn)
			}()
		default:
			// We can't acquire a semaphore slot, so log the rejection
			log.Printf("%s Connection rejected: maximum concurrent connections reached", config.ServerName)
		}
	}
}
