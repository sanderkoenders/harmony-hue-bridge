package httpapi

import (
	"log"
	"net/http"

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

func (s *Server) Run(addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/description.xml", handlers.HandleDescription(s.logger, s.bridge))
	mux.HandleFunc("/api/", handlers.HandleApi(s.logger, s.bridge))
	mux.HandleFunc("/api/001788FFFE23BFC2/lights", handlers.HandleLights(s.logger))
	mux.HandleFunc("/api/001788FFFE23BFC2/lights/1", handlers.HandleLight(s.logger))
	mux.HandleFunc("/api/001788FFFE23BFC2/groups", handlers.HandleGroups(s.logger))
	mux.HandleFunc("/api/001788FFFE23BFC2/scenes", handlers.HandleScenes(s.logger))

	mux.HandleFunc("/", s.HandleUnknownRequest)

	s.logger.Printf("HTTP server listening on %s", addr)

	return http.ListenAndServe(addr, mux)
}
