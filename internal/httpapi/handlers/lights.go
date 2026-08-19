package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/sanderkoenders/harmony-hue-bridge/internal/mqtt"
)

func HandleLights(logger *log.Logger, m mqtt.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("HTTP %s %s from %s", r.Method, r.URL.RequestURI(), r.RemoteAddr)

		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		lights, err := m.GetLights(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Build hue-like response mapping id -> details
		resp := make(map[string]any)
		for id, l := range lights {
			resp[id] = map[string]any{
				"state": map[string]any{
					"on":  l.On,
					"bri": l.Brightness,
					"hue": l.Hue,
					"sat": l.Sat,
				},
				"type":             "Extended color light",
				"name":             l.Name,
				"modelid":          "LCT001",
				"manufacturername": "Signify",
				"swversion":        "1.0",
			}
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		_ = enc.Encode(resp)
	}
}

func HandleLight(logger *log.Logger, m mqtt.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s", r.Method, r.URL.Path)

		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/api/001788FFFE23BFC2/lights/")

		lights, err := m.GetLights(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		l, ok := lights[id]
		if !ok {
			http.NotFound(w, r)
			return
		}

		resp := map[string]any{
			"state": map[string]any{
				"on":  l.On,
				"bri": l.Brightness,
				"hue": l.Hue,
				"sat": l.Sat,
			},
			"type":             "Extended color light",
			"name":             l.Name,
			"modelid":          "LCT001",
			"manufacturername": "HarmonyHueBridge",
			"uniqueid":         fmt.Sprintf("00:11:22:33:44:55:66:%s-00", id),
			"swversion":        "1.0",
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		_ = enc.Encode(resp)
	}
}
