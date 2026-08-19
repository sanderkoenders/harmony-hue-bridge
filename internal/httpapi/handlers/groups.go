package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetGroups(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{
				"1": {
					"name": "Living Room",
					"lights": ["1", "2"],
					"type": "Room",
					"class": "Living room"
				}
			}`)
			return
		case http.MethodPost:
			body, _ := io.ReadAll(r.Body)
			logger.Printf("GROUPS POST body: %s", body)

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"success":{"id":"1"}}`)
			return
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}
