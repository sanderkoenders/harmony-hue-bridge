package ssdp

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/sanderkoenders/harmony-hue-bridge/internal/bridge"
	"github.com/sanderkoenders/harmony-hue-bridge/internal/ssdp/handlers"
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
	conn, err := s.listen(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	return s.readLoop(ctx, conn)
}

func (s *Server) listen(ctx context.Context) (*net.UDPConn, error) {
	addr := &net.UDPAddr{
		IP:   net.ParseIP(MulticastAddress),
		Port: Port,
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("listen on SSDP multicast: %w", err)
	}

	s.logger.Printf(
		"SSDP listening on %s:%d",
		MulticastAddress,
		Port,
	)

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	return conn, nil
}

func (s *Server) readLoop(ctx context.Context, conn *net.UDPConn) error {
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
	msg, err := FromDataFrame(data)
	if err != nil {
		s.logger.Printf("invalid SSDP message: %v", err)
		return
	}

	if msg.IsMSearch() {
		handlers.HandleMSearch(s.logger, s.bridge, conn, remoteAddr, msg.Headers)
		return
	}

	s.logger.Printf("Unhandled SSDP message: method=%q headers=%v", msg.Method, msg.Headers)
}
