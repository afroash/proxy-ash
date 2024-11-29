package api

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"

	"github.com/afroash/proxy-ash/internal/config"
	"github.com/afroash/proxy-ash/internal/metrics"
)

type Server struct {
	cfg     *config.Config
	metrics *metrics.Collector
}

func NewServer(cfg *config.Config, metrics *metrics.Collector) *Server {
	return &Server{
		cfg:     cfg,
		metrics: metrics,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp4", s.cfg.ListenAddr)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer listener.Close()

	fmt.Printf("Proxy server listening on %s\n", s.cfg.ListenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(downstream net.Conn) {
	connID := fmt.Sprintf("%s-%d", downstream.RemoteAddr(), time.Now().UnixNano())
	s.metrics.StartConnection(connID)
	defer func() {
		s.metrics.EndConnection(connID)
		downstream.Close()
	}()

	upstream, err := net.Dial("tcp", s.cfg.UpstreamAddr)
	if err != nil {
		return
	}
	defer upstream.Close()

	done := make(chan bool, 2)
	go s.pipe(upstream, downstream, done, connID)
	go s.pipe(downstream, upstream, done, connID)

	<-done
}

func (s *Server) pipe(dst io.Writer, src io.Reader, done chan bool, connID string) {
	_ = connID
	buffer := make([]byte, 32*1024)
	for {
		// Simulate packet loss if enabled
		if s.cfg.PacketLoss.Enabled && rand.Float64() < s.cfg.PacketLoss.Percentage {
			s.metrics.RecordPacketLoss()
			continue
		}

		// Simulate latency if enabled
		if s.cfg.Latency.Enabled {
			time.Sleep(s.cfg.Latency.Duration)
		}

		n, err := src.Read(buffer)
		if err != nil {
			break
		}

		_, err = dst.Write(buffer[:n])
		if err != nil {
			break
		}

		s.metrics.RecordBytes(n)
	}
	done <- true
}
