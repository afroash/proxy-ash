package config

import "time"

type Config struct {
	ListenAddr   string
	UpstreamAddr string

	// Network condition simulation settings
	Latency struct {
		Enabled  bool          `yaml:"enabled"`
		Min      int           `yaml:"min_ms"`
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
	// TODO lmplematation for loading config from yaml file. maybe we use json.
	// not sure.
	return &Config{}, nil
}
