package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sanderkoenders/harmony-hue-bridge/internal/ssdp"
	"github.com/sanderkoenders/harmony-hue-bridge/internal/httpapi"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags | log.Lmicroseconds)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	ssdpServer := ssdp.NewServer(logger)

	httpServer := httpapi.NewServer(logger)

	go func() {
		if err := httpServer.Run(":80"); err != nil {
			logger.Printf("HTTP server stopped: %v", err)
		}
	}()

	if err := ssdpServer.Run(ctx); err != nil {
		logger.Fatal(err)
	}
}
