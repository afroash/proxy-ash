package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/afroash/proxy-ash/api"
	"github.com/afroash/proxy-ash/internal/config"
	"github.com/afroash/proxy-ash/internal/metrics"
	"github.com/afroash/proxy-ash/internal/proxy"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Check config has loaded
	if cfg.ListenAddr == "" || cfg.UpstreamAddr == "" {
		log.Fatal("Missing required configuration parameters")
	}

	// Initialize metrics collector
	collector := metrics.NewCollector()

	// Create servers
	proxyServer := proxy.NewServer(cfg, collector)
	apiServer := api.NewServer(cfg, collector)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Use WaitGroup to manage both servers
	var wg sync.WaitGroup
	wg.Add(2)

	// Start proxy server
	go func() {
		defer wg.Done()
		if err := proxyServer.Start(); err != nil {
			log.Printf("Proxy server error: %v", err)
			sigChan <- syscall.SIGTERM
		}
	}()

	// Start API server
	go func() {
		defer wg.Done()
		if err := apiServer.Start(); err != nil {
			log.Printf("API server error: %v", err)
			sigChan <- syscall.SIGTERM
		}
	}()

	// Wait for shutdown signal
	sig := <-sigChan
	log.Printf("Received signal %v, shutting down...", sig)

	// Shutdown both servers
	if err := proxyServer.Shutdown(); err != nil {
		log.Printf("Error shutting down proxy server: %v", err)
	}
	if err := apiServer.Shutdown(); err != nil {
		log.Printf("Error shutting down API server: %v", err)
	}

	// Wait for both servers to finish
	wg.Wait()
	log.Println("Shutdown complete")
}
