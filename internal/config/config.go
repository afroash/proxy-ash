package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ListenAddr   string `yaml:"ListenAddr"`
	UpstreamAddr string `yaml:"UpstreamAddr"`
	//Metrics      MetricsConfig `yaml:"metrics"`

	// Network condition simulation settings
	Latency struct {
		Enabled  bool          `yaml:"Enabled"`
		Min      int           `yaml:"Min"`
		Max      int           `yaml:"max_ms"`
		Duration time.Duration `yaml:"duration"`
	} `yaml:"latency"`

	PacketLoss struct {
		Enabled    bool    `yaml:"enabled"`
		Percentage float64 `yaml:"percentage"`
	} `yaml:"packet_loss"`

	Bandwidth struct {
		Enabled bool `yaml:"enabled"`
		Limit   int  `yaml:"limit_kbps"`
	} `yaml:"bandwidth"`

	// Metrics collection settings
	Metrics struct {
		Enabled bool   `yaml:"enabled"`
		Path    string `yaml:"path"`
	} `yaml:"metrics"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
