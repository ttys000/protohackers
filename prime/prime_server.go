package prime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"protohackers/servers"
)

const serverName = "[Prime]"

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
	decoder := json.NewDecoder(conn)
	req := &isPrimeRequest{}
	if err := decoder.Decode(req); err != nil {
		log.Printf("%s JSON decode error: %v", serverName, err)
		io.WriteString(conn, "malformed response")
		return
	}

	if req.Method != "isPrime" {
		log.Printf("%s The value of method field is not isPrime: %v", serverName, req.Method)
		io.WriteString(conn, "malformed response")
		return
	}

	resp := &isPrimeResonse{
		Method: "isPrime",
		Prime:  isPrime(req.Number),
	}

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(resp); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

func isPrime(n float64) bool {
	// Check if number is an integer
	if n != float64(int(n)) {
		return false
	}

	// Convert to int for easier calculation
	num := int(n)

	// Numbers less than 2 are not prime
	if num < 2 {
		return false
	}

	// Check divisibility up to square root
	for i := 2; i*i <= num; i++ {
		if num%i == 0 {
			return false
		}
	}

	return true
}

type isPrimeRequest struct {
	Method string  `json:"method"`
	Number float64 `json:"number"`
}

type isPrimeResonse struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}
