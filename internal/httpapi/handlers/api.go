package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sanderkoenders/harmony-hue-bridge/internal/bridge"
)

func HandleApi(logger *log.Logger, bridge *bridge.Bridge) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("HTTP %s %s from %s", r.Method, r.URL.RequestURI(), r.RemoteAddr)

		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `[]`)
			return
		case http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `[
				{
					"success": {
						"username": `+bridge.ID+`
					}
				}
			]`)
			return
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}
