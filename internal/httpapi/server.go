package httpapi

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Light struct {
	ID   string
	Name string
}

var lights = []Light{
	{ID: "1", Name: "Living Room"},
	{ID: "2", Name: "Kitchen"},
	{ID: "3", Name: "Bedroom"},
}

const descriptionXML = `<?xml version="1.0" encoding="UTF-8" ?>
<root xmlns="urn:schemas-upnp-org:device-1-0">
	<specVersion>
		<major>1</major>
		<minor>0</minor>
	</specVersion>

	<URLBase>http://192.168.30.104:80/</URLBase>

	<device>
		<deviceType>urn:schemas-upnp-org:device:Basic:1</deviceType>

		<friendlyName>Home Assistant Bridge (192.168.30.104)</friendlyName>

		<manufacturer>Royal Philips Electronics</manufacturer>
		<manufacturerURL>http://www.philips.com</manufacturerURL>

		<modelDescription>Philips hue Personal Wireless Lighting</modelDescription>
		<modelName>Philips hue bridge 2015</modelName>
		<modelNumber>BSB002</modelNumber>
		<modelURL>http://www.meethue.com</modelURL>

		<serialNumber>001788FFFE23BFC2</serialNumber>

		<UDN>uuid:2f402f80-da50-11e1-9b23-001788255acc</UDN>
	</device>
</root>`

type Server struct {
	logger *log.Logger
}

func NewServer(logger *log.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func handleAPI(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("HTTP %s %s from %s", r.Method, r.URL.RequestURI(), r.RemoteAddr)

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, `[
			{
				"success": {
					"username": "001788FFFE23BFC2"
				}
			}
		]`)
	}
}

func handleLights(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s", r.Method, r.URL.Path)

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

func handleLight(logger *log.Logger) http.HandlerFunc {
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

func handleGroups(logger *log.Logger) http.HandlerFunc {
	// return func(w http.ResponseWriter, r *http.Request) {
	// 	start := time.Now()

	// 	logger.Printf("GET %s", r.URL.Path)

	// 	if r.Method != http.MethodGet {
	// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 		return
	// 	}

	// 	w.Header().Set("Content-Type", "application/json")

	// 	fmt.Fprint(w, `{
	// 		"1": {
	// 			"name": "Living Room",
	// 			"lights": ["1", "2"],
	// 			"type": "Room",
	// 			"class": "Living room"
	// 		}
	// 	}`)

	// 	logger.Printf("GET /groups completed in %s", time.Since(start))
	// }

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		body, _ := io.ReadAll(r.Body)
		logger.Printf("GROUPS POST body: %s", body)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"success":{"id":"1"}}`)
	}
}

func handleScenes(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s", r.Method, r.URL.Path)

		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, `{}`)
	}
}

func (s *Server) Run(addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/description.xml", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		s.logger.Printf(
			"HTTP request: %s %s from %s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		)

		w.Header().Set(
			"Content-Type",
			"text/xml; charset=\"utf-8\"",
		)

		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(descriptionXML))
	})

	mux.HandleFunc("/api/", handleAPI(s.logger))

	mux.HandleFunc("/api/001788FFFE23BFC2/lights", handleLights(s.logger))
	mux.HandleFunc("/api/001788FFFE23BFC2/lights/1", handleLight(s.logger))
	mux.HandleFunc("/api/001788FFFE23BFC2/groups", handleGroups(s.logger))
	mux.HandleFunc("/api/001788FFFE23BFC2/scenes", handleScenes(s.logger))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Printf("HTTP %s %s from %s", r.Method, r.URL.RequestURI(), r.RemoteAddr)

		for name, values := range r.Header {
			for _, value := range values {
				s.logger.Printf("  %s: %s", name, value)
			}
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.Printf("  Error reading body: %v", err)
		} else {
			s.logger.Printf("  Body: %q", string(body))
		}

		http.NotFound(w, r)
	})

	s.logger.Printf("HTTP server listening on %s", addr)

	return http.ListenAndServe(addr, mux)
}
