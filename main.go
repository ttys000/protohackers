package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"protohackers/echo"
	"protohackers/price"
	"protohackers/prime"
)

func main() {
	// Create a context that will be canceled on interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Create a wait group to track active connections
	var wg sync.WaitGroup

	// Start servers
	go echo.RunServer(ctx, 5001, &wg)
	go prime.RunServer(ctx, 5002, &wg)
	go price.RunServer(ctx, 5003, &wg)

	// Wait for interrupt signal
	<-ctx.Done()
	log.Println("Shutdown signal received, waiting for active connections to complete...")

	// Wait for all active connections to finish
	wg.Wait()
	log.Println("All servers shutdown complete")
}
