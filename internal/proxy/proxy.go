package proxy

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/afroash/proxy-ash/internal/config"
	"github.com/afroash/proxy-ash/internal/metrics"
)

type Server struct {
	cfg      *config.Config
	metrics  *metrics.Collector
	listener net.Listener
	wg       sync.WaitGroup
	shutdown chan struct{}
}

func NewServer(cfg *config.Config, metrics *metrics.Collector) *Server {
	return &Server{
		cfg:      cfg,
		metrics:  metrics,
		shutdown: make(chan struct{}),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.cfg.ListenAddr)
	if err != nil {
		return fmt.Errorf("failed to start proxy listener: %w", err)
	}
	s.listener = listener

	fmt.Printf("Proxy server listening on %s\n", s.cfg.ListenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-s.shutdown:
				return nil
			default:
				fmt.Printf("Failed to accept connection: %v\n", err)
				continue
			}
		}

		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(downstream net.Conn) {
	defer s.wg.Done()

	connID := fmt.Sprintf("%s-%d", downstream.RemoteAddr(), time.Now().UnixNano())
	s.metrics.StartConnection(connID)
	defer func() {
		s.metrics.EndConnection(connID)
		downstream.Close()
	}()

	upstream, err := net.Dial("tcp", s.cfg.UpstreamAddr)
	if err != nil {
		fmt.Printf("Failed to connect to upstream %s: %v\n", s.cfg.UpstreamAddr, err)
		return
	}
	fmt.Printf("Connected to upstream %s\n", s.cfg.UpstreamAddr)
	defer upstream.Close()

	// Create channels for signaling completion
	done := make(chan struct{})

	// Start bidirectional copying
	go func() {
		s.pipe(upstream, downstream, connID)
		close(done)
	}()

	s.pipe(downstream, upstream, connID)
	<-done
}

func (s *Server) pipe(dst io.Writer, src io.Reader, connID string) {
	buffer := make([]byte, 32*1024)
	for {
		// Simulate packet loss if enabled
		if s.cfg.PacketLoss.Enabled && rand.Float64() < s.cfg.PacketLoss.Percentage/100.0 {
			s.metrics.RecordPacketLoss()
			continue
		}

		// Simulate latency if enabled
		if s.cfg.Latency.Enabled {
			time.Sleep(s.cfg.Latency.Duration)
		}

		n, err := src.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading from connection %s: %v\n", connID, err)
			}
			return
		}

		// Simulate bandwidth limit if enabled
		if s.cfg.Bandwidth.Enabled && s.cfg.Bandwidth.Limit > 0 {
			delay := time.Duration(float64(n) * 8 / float64(s.cfg.Bandwidth.Limit) * float64(time.Second))
			time.Sleep(delay)
		}

		_, err = dst.Write(buffer[:n])
		if err != nil {
			fmt.Printf("Error writing to connection %s: %v\n", connID, err)
			return
		}

		s.metrics.RecordBytes(n)
	}
}

func (s *Server) Shutdown() error {
	close(s.shutdown)
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return fmt.Errorf("error closing listener: %w", err)
		}
	}
	s.wg.Wait()
	return nil
}
