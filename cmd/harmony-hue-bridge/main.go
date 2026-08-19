package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sanderkoenders/harmony-hue-bridge/internal/bridge"
	"github.com/sanderkoenders/harmony-hue-bridge/internal/httpapi"
	"github.com/sanderkoenders/harmony-hue-bridge/internal/ssdp"
)

func main() {
	hueBridge := bridge.New(
		"001788FFFE23BFC2",
		"2f402f80-da50-11e1-9b23-001788255acc",
		"192.168.1.104",
		8080,
	)

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	ctx, stop := newContext()
	defer stop()

	ssdpServer := ssdp.NewServer(logger, hueBridge)
	httpServer := httpapi.NewServer(logger, hueBridge)

	httpErrCh := make(chan error, 1)
	go func() {
		httpErrCh <- httpServer.Run(ctx, ":"+fmt.Sprint(hueBridge.Port))
	}()

	if err := ssdpServer.Run(ctx); err != nil {
		logger.Fatal(err)
	}

	logger.Printf("SSDP server shutdown successful")

	if err := <-httpErrCh; err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Printf("HTTP server stopped: %v", err)
		}
		return
	}

	logger.Printf("HTTP server shutdown successful")
}

func newContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
}
