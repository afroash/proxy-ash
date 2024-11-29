package proxy

import (
	"fmt"
	"net"
	"sync"

	"github.com/afroash/proxy-ash/internal/config"
	"github.com/afroash/proxy-ash/internal/simulator"
)

type Server struct {
	cfg      *config.Config
	listener net.Listener
	wg       sync.WaitGroup
	shutdown chan struct{}
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:      cfg,
		shutdown: make(chan struct{}),
		//metrics:  metrics.NewCollector(),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.cfg.ListenAddr)
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}
	s.listener = listener

	fmt.Printf("Proxy listening on %s\n", s.cfg.ListenAddr)

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
	defer downstream.Close()

	// Connect to upstream
	upstream, err := net.Dial("tcp", s.cfg.UpstreamAddr)
	if err != nil {
		fmt.Printf("Failed to connect to upstream: %v\n", err)
		return
	}
	defer upstream.Close()

	// Create network condition simulator
	sim := simulator.NewSimulator(s.cfg)

	// Create bidirectional channels
	upLink := NewLink(sim)
	downLink := NewLink(sim)

	// Start proxying
	done := make(chan bool, 2)
	go proxyTraffic(downstream, upLink, done)
	go proxyTraffic(upstream, downLink, done)

	<-done
	<-done
}

func (s *Server) Shutdown() {
	close(s.shutdown)
	if s.listener != nil {
		s.listener.Close()
	}
	s.wg.Wait()
}
