package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

type Collector struct {
	stats       Stats
	connections map[string]*ConnectionStats
	mu          sync.RWMutex
}

func NewCollector() *Collector {
	c := &Collector{
		connections: make(map[string]*ConnectionStats),
	}
	go c.startReporter()
	return c
}

func (c *Collector) RecordBytes(n int) {
	atomic.AddInt64(&c.stats.TotalBytes, int64(n))
}

func (c *Collector) RecordLatency(latency time.Duration) {
	// Use exponential moving average for latency
	const alpha = 0.1
	current := c.stats.AverageLatency
	new := current*(1-alpha) + float64(latency.Milliseconds())*alpha
	c.stats.AverageLatency = new
}

func (c *Collector) RecordPacketLoss() {
	return
}

func (c *Collector) StartConnection(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	atomic.AddInt64(&c.stats.ActiveConnections, 1)
	atomic.AddInt64(&c.stats.TotalConnections, 1)

	c.connections[id] = &ConnectionStats{
		StartTime: time.Now().UnixNano(),
	}
}

func (c *Collector) EndConnection(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	atomic.AddInt64(&c.stats.ActiveConnections, -1)

	if conn, exists := c.connections[id]; exists {
		conn.EndTime = time.Now().UnixNano()
		// Could store historical data here
		delete(c.connections, id)
	}
}

func (c *Collector) GetStats() Stats {
	return Stats{
		ActiveConnections: atomic.LoadInt64(&c.stats.ActiveConnections),
		TotalConnections:  atomic.LoadInt64(&c.stats.TotalConnections),
		TotalBytes:        atomic.LoadInt64(&c.stats.TotalBytes),
		AverageLatency:    c.stats.AverageLatency,
	}
}

func (c *Collector) startReporter() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats := c.GetStats()
		// Here you could:
		// 1. Export to Prometheus
		// 2. Write to a time-series database
		// 3. Send to a monitoring service
		// 4. Log to file
		// For now, we'll just print
		println("Active Connections:", stats.ActiveConnections)
		println("Total Bytes:", stats.TotalBytes)
		println("Average Latency:", stats.AverageLatency)
	}
}
