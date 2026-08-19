package ssdp

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
)

const descriptionXML = `<?xml version="1.0" encoding="UTF-8" ?>
<root xmlns="urn:schemas-upnp-org:device-1-0">
	<specVersion>
		<major>1</major>
		<minor>0</minor>
	</specVersion>

	<URLBase>http://192.168.20.203:80/</URLBase>

	<device>
		<deviceType>urn:schemas-upnp-org:device:Basic:1</deviceType>

		<friendlyName>Home Assistant Bridge (192.168.20.203)</friendlyName>

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

const (
	MulticastAddress = "239.255.255.250"
	Port             = 1900
)

type Server struct {
	logger *log.Logger
}

func NewServer(logger *log.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	addr := &net.UDPAddr{
		IP:   net.ParseIP(MulticastAddress),
		Port: Port,
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return fmt.Errorf("listen on SSDP multicast: %w", err)
	}
	defer conn.Close()

	s.logger.Printf(
		"SSDP listening on %s:%d",
		MulticastAddress,
		Port,
	)

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	buf := make([]byte, 64*1024)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				return fmt.Errorf("read SSDP packet: %w", err)
			}
		}

		s.handlePacket(conn, remoteAddr, buf[:n])
	}
}

func (s *Server) handlePacket(
	conn *net.UDPConn,
	remoteAddr *net.UDPAddr,
	data []byte,
) {
	msg, err := ParseMessage(data)
	if err != nil {
		s.logger.Printf("invalid SSDP message: %v", err)
		return
	}

	if !msg.IsMSearch() {
		return
	}

	st := msg.Header("st")
	man := msg.Header("man")

	// s.logger.Printf(
	// 	"M-SEARCH from %s: ST=%q MAN=%q",
	// 	remoteAddr,
	// 	st,
	// 	man,
	// )

	// Only respond to actual SSDP discovery requests.
	if !strings.EqualFold(man, `"ssdp:discover"`) {
		return
	}

	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n"+
			"CACHE-CONTROL: max-age=100\r\n"+
			"EXT:\r\n"+
			"LOCATION: http://192.168.30.104:80/description.xml\r\n"+
			"SERVER: Linux/1.0 UPnP/1.0 IpBridge/1.0\r\n"+
			"ST: %s\r\n"+
			"USN: uuid:2f402f80-da50-11e1-9b23-001788255acc::urn:schemas-upnp-org:device:basic:1\r\n"+
			"\r\n",
		st,
	)

	_, err = conn.WriteToUDP([]byte(response), remoteAddr)
	if err != nil {
		s.logger.Printf(
			"failed to send SSDP response to %s: %v",
			remoteAddr,
			err,
		)
		return
	}

	// s.logger.Printf(
	// 	"SSDP response sent to %s: ST=%q",
	// 	remoteAddr,
	// 	st,
	// )
}
