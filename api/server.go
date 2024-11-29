package api

import (
	"fmt"
	"net/http"

	"github.com/afroash/proxy-ash/internal/config"
	"github.com/afroash/proxy-ash/internal/metrics"
)

type Server struct {
	cfg     *config.Config
	metrics *metrics.Collector
	server  *http.Server
}

func NewServer(cfg *config.Config, metrics *metrics.Collector) *Server {
	return &Server{
		cfg:     cfg,
		metrics: metrics,
	}
}

func (s *Server) Start() error {
	// Set up HTTP routes
	mux := http.NewServeMux()

	// Metrics endpoint
	mux.HandleFunc("/metrics", s.handleMetrics)

	// Create HTTP server
	s.server = &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", 8081), // Use a different port for the API server
		Handler: mux,
	}

	fmt.Printf("API server listening on %s\n", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if !s.cfg.Metrics.Enabled {
		http.Error(w, "Metrics are disabled", http.StatusServiceUnavailable)
		return
	}

	// Get metrics data
	stats := s.metrics.GetStats()

	// Write metrics response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
		"active_connections": %d,
		"total_bytes": %d,
		"average_latency": %f
	}`, stats.ActiveConnections, stats.TotalBytes, stats.AverageLatency)
}

func (s *Server) Shutdown() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}
