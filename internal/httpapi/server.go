package httpapi

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/sanderkoenders/harmony-hue-bridge/internal/bridge"
	"github.com/sanderkoenders/harmony-hue-bridge/internal/httpapi/handlers"
)

type Server struct {
	logger *log.Logger
	bridge *bridge.Bridge
}

func NewServer(logger *log.Logger, bridge *bridge.Bridge) *Server {
	return &Server{
		logger: logger,
		bridge: bridge,
	}
}

func (s *Server) HandleUnknownRequest(w http.ResponseWriter, r *http.Request) {
	s.logger.Printf("Unknown HTTP %s %s from %s", r.Method, r.URL.RequestURI(), r.RemoteAddr)

	http.NotFound(w, r)
}

func (s *Server) Run(ctx context.Context, addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/description.xml", handlers.HandleDescription(s.logger, s.bridge))
	mux.HandleFunc("/api/", handlers.HandleApi(s.logger, s.bridge))
	mux.HandleFunc("/api/"+s.bridge.Username+"/lights", handlers.HandleLights(s.logger))
	mux.HandleFunc("/api/"+s.bridge.Username+"/lights/1", handlers.HandleLight(s.logger))
	mux.HandleFunc("/api/"+s.bridge.Username+"/groups", handlers.HandleGroups(s.logger))
	mux.HandleFunc("/api/"+s.bridge.Username+"/scenes", handlers.HandleScenes(s.logger))
	mux.HandleFunc("/", s.HandleUnknownRequest)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			s.logger.Printf("HTTP shutdown failed: %v", err)
		}
	}()

	s.logger.Printf("HTTP server listening on %s", addr)

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
