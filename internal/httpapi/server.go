package httpapi

import (
	"io"
	"log"
	"net/http"

	"github.com/sanderkoenders/harmony-hue-bridge/internal/bridge"
	"github.com/sanderkoenders/harmony-hue-bridge/internal/httpapi/handlers"
)

type Light struct {
	ID   string
	Name string
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
	bridge *bridge.Bridge
}

func NewServer(logger *log.Logger, bridge *bridge.Bridge) *Server {
	return &Server{
		logger: logger,
		bridge: bridge,
	}
}

func (s *Server) Run(addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/description.xml", handlers.GetDescription(s.logger))

	mux.HandleFunc("/api/", handlers.GetAPI(s.logger, s.bridge))
	mux.HandleFunc("/api/001788FFFE23BFC2/lights", handlers.GetLights(s.logger))
	mux.HandleFunc("/api/001788FFFE23BFC2/lights/1", handlers.GetLight(s.logger))
	mux.HandleFunc("/api/001788FFFE23BFC2/groups", handlers.GetGroups(s.logger))
	mux.HandleFunc("/api/001788FFFE23BFC2/scenes", handlers.GetScenes(s.logger))

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
