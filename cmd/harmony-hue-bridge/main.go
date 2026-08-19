package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sanderkoenders/harmony-hue-bridge/internal/bridge"
	"github.com/sanderkoenders/harmony-hue-bridge/internal/httpapi"
	"github.com/sanderkoenders/harmony-hue-bridge/internal/ssdp"
)

func main() {
	bridge := bridge.New(
		"001788FFFE23BFC2",
		"2f402f80-da50-11e1-9b23-001788255acc",
		"192.168.1.104",
		8080,
	)

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	ssdpServer := ssdp.NewServer(logger, bridge)

	httpServer := httpapi.NewServer(logger, bridge)

	go func() {
		if err := httpServer.Run(":" + fmt.Sprint(bridge.Port)); err != nil {
			logger.Printf("HTTP server stopped: %v", err)
		}
	}()

	if err := ssdpServer.Run(ctx); err != nil {
		logger.Fatal(err)
	}
}
