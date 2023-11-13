package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func GracefuleShutdown() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(
		sigChannel,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-sigChannel
		// Handle Shutdown
		fmt.Println("Shutdown signal has been captured")
		os.Exit(0)
	}()
}
