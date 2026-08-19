package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func HandleLights(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("HTTP %s %s from %s", r.Method, r.URL.RequestURI(), r.RemoteAddr)

		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, `{
			"1": {
				"state": {
					"on": false,
					"bri": 254,
					"hue": 10000,
					"sat": 200,
					"effect": "none",
					"xy": [0.5, 0.5],
					"ct": 250,
					"alert": "none",
					"colormode": "ct"
				},
				"type": "Extended color light",
				"name": "Living Room 1",
				"modelid": "LCT001",
				"manufacturername": "Signify",
				"productname": "Hue color lamp",
				"uniqueid": "00:11:22:33:44:55:66:01-0b",
				"swversion": "1.0"
			},
			"2": {
				"state": {
					"on": false,
					"bri": 254,
					"hue": 20000,
					"sat": 200,
					"effect": "none",
					"xy": [0.6, 0.4],
					"ct": 300,
					"alert": "none",
					"colormode": "ct"
				},
				"type": "Extended color light",
				"name": "Living Room 2",
				"modelid": "LCT001",
				"manufacturername": "Signify",
				"productname": "Hue color lamp",
				"uniqueid": "00:11:22:33:44:55:66:02-0b",
				"swversion": "1.0"
			},
			"3": {
				"state": {
					"on": false,
					"bri": 254,
					"hue": 20000,
					"sat": 200,
					"effect": "none",
					"xy": [0.6, 0.4],
					"ct": 300,
					"alert": "none",
					"colormode": "ct"
				},
				"type": "Extended color light",
				"name": "Living Room 3",
				"modelid": "LCT001",
				"manufacturername": "Signify",
				"productname": "Hue color lamp",
				"uniqueid": "00:11:22:33:44:55:66:03-0b",
				"swversion": "1.0"
			}
		}`)
	}
}

func HandleLight(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s", r.Method, r.URL.Path)

		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/api/001788FFFE23BFC2/lights/")

		var name string

		switch id {
		case "1":
			name = "Living Room"
		case "2":
			name = "Kitchen"
		case "3":
			name = "Bedroom"
		default:
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprintf(w, `{
			"state": {
				"on": false,
				"bri": 254,
				"hue": 0,
				"sat": 0,
				"effect": "none",
				"xy": [0.5, 0.5],
				"ct": 250,
				"alert": "none",
				"colormode": "ct"
			},
			"type": "Extended color light",
			"name": %q,
			"modelid": "LCT001",
			"manufacturername": "HarmonyHueBridge",
			"uniqueid": "00:11:22:33:44:55:66:%s-00",
			"swversion": "1.0"
		}`, name, id)
	}
}
