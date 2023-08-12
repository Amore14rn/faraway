package utils

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// GracefulShutdown provides a way to initiate a graceful shutdown
func GracefulShutdown(ctx context.Context) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		fmt.Printf("Received signal: %s\n", sig)
	case <-ctx.Done():
		fmt.Println("Graceful shutdown initiated")
	}

	// Perform cleanup tasks if necessary before exiting
}
