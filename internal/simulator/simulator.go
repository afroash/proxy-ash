package simulator

import (
	"time"

	"github.com/afroash/proxy-ash/internal/config"
	"golang.org/x/exp/rand"
)

type Simulator struct {
	cfg *config.Config
}

func NewSimulator(cfg *config.Config) *Simulator {
	return &Simulator{
		cfg: cfg,
	}
}

// ApplyConditions applies all configured network conditions to the data
func (s *Simulator) ApplyConditions(data []byte) {
	s.ApplyLatency()
	
	if s.ShouldDropPacket() {
		data = nil
		return
	}

	if s.cfg.Bandwidth.Enabled {
		s.ApplyBandwidthLimit(len(data))
	}
}

func (s *Simulator) ApplyLatency() {
	if !s.cfg.Latency.Enabled {
		return
	}

	min := s.cfg.Latency.Min
	max := s.cfg.Latency.Max

	if min == max {
		time.Sleep(time.Duration(min) * time.Millisecond)
		return
	}

	// Random latency between min and max
	latency := rand.Intn(max-min) + min
	time.Sleep(time.Duration(latency) * time.Millisecond)
}

func (s *Simulator) ShouldDropPacket() bool {
	if !s.cfg.PacketLoss.Enabled {
		return false
	}

	return rand.Float64() < s.cfg.PacketLoss.Percentage
}

func (s *Simulator) ApplyBandwidthLimit(bytes int) {
	if !s.cfg.Bandwidth.Enabled {
		return
	}

	// Calculate delay based on bandwidth limit
	bytesPerSecond := s.cfg.Bandwidth.Limit * 1024 / 8 // Convert kbps to bytes per second
	delay := time.Duration(float64(bytes) / float64(bytesPerSecond) * float64(time.Second))
	time.Sleep(delay)
}
