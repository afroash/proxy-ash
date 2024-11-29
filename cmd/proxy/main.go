package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/afroash/proxy-ash/api"
	"github.com/afroash/proxy-ash/internal/config"
	"github.com/afroash/proxy-ash/internal/metrics"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	//check config has loaded.
	if cfg.ListenAddr == "" || cfg.UpstreamAddr == "" {
		log.Println("Missing required configuration parameters")
	}

	// Initialize metrics collector
	collector := metrics.NewCollector()

	// Create and start server
	server := api.NewServer(cfg, collector)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Server error: %v", err)
			sigChan <- syscall.SIGTERM
		}
	}()

	// Wait for shutdown signal
	sig := <-sigChan
	log.Printf("Received signal %v, shutting down...", sig)
}
