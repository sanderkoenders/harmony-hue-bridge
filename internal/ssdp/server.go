package ssdp

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/sanderkoenders/harmony-hue-bridge/internal/bridge"
)

const (
	MulticastAddress = "239.255.255.250"
	Port             = 1900
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

	// Only respond to actual SSDP discovery requests.
	if !strings.EqualFold(man, `"ssdp:discover"`) {
		return
	}

	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n"+
			"CACHE-CONTROL: max-age=100\r\n"+
			"EXT:\r\n"+
			"LOCATION: http://"+s.bridge.IpAddr+":"+fmt.Sprintf("%d", s.bridge.Port)+"/description.xml\r\n"+
			"SERVER: Linux/1.0 UPnP/1.0 IpBridge/1.0\r\n"+
			"ST: %s\r\n"+
			"USN: uuid:"+s.bridge.UUID+"::urn:schemas-upnp-org:device:basic:1\r\n"+
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
}
