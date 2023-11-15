package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Closable interface {
	GracefulShutdown()
}

func GracefulShutdown(c []Closable) {
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
		fmt.Println("Closing Closables")

		for _, closable := range c {
			closable.GracefulShutdown()
		}
		fmt.Println("Closed all Closables")

		os.Exit(0)
	}()
}
