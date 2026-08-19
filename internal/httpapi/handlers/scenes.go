package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func HandleScenes(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("HTTP %s %s from %s", r.Method, r.URL.RequestURI(), r.RemoteAddr)

		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, `{}`)
	}
}
