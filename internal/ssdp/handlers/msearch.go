package handlers

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/sanderkoenders/harmony-hue-bridge/internal/bridge"
)

func HandleMSearch(logger *log.Logger, bridge *bridge.Bridge, conn *net.UDPConn, remoteAddr *net.UDPAddr, headers map[string]string) {
	st := headers["st"]
	man := headers["man"]

	if !strings.EqualFold(man, `"ssdp:discover"`) {
		logger.Printf("ignored SSDP M-SEARCH message without ssdp:discover MAN header")
		return
	}

	response := buildMSearchResponse(st, bridge)

	_, err := conn.WriteToUDP([]byte(response), remoteAddr)
	if err != nil {
		logger.Printf("failed to send SSDP response to %s: %v", remoteAddr, err)
	}
}

func buildMSearchResponse(st string, bridge *bridge.Bridge) string {
	location := fmt.Sprintf("http://%s:%d/description.xml", bridge.IpAddr, bridge.Port)
	usn := fmt.Sprintf("uuid:%s::upnp:rootdevice", bridge.UUID)

	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n"+
			"CACHE-CONTROL: max-age=100\r\n"+
			"EXT: \r\n"+
			"LOCATION: %s\r\n"+
			"SERVER: Linux/1.0 UPnP/1.0 IpBridge/1.0\r\n"+
			"ST: %s\r\n"+
			"USN: %s\r\n\r\n",
		location,
		st,
		usn,
	)

	return response
}
