package proxy

import (
	"io"

	"github.com/afroash/proxy-ash/internal/simulator"
)

// Link represents a network connection with simulated conditions
type Link struct {
	sim *simulator.Simulator
	ch  chan []byte
}

// NewLink creates a new Link with the given simulator
func NewLink(sim *simulator.Simulator) *Link {
	return &Link{
		sim: sim,
		ch:  make(chan []byte, 1024),
	}
}

// Read implements io.Reader interface
func (l *Link) Read(p []byte) (n int, err error) {
	data, ok := <-l.ch
	if !ok {
		return 0, io.EOF
	}

	// Apply network conditions
	l.sim.ApplyConditions(data)

	n = copy(p, data)
	return n, nil
}

// Write implements io.Writer interface
func (l *Link) Write(p []byte) (n int, err error) {
	data := make([]byte, len(p))
	copy(data, p)
	l.ch <- data
	return len(p), nil
}

// proxyTraffic handles the data transfer between connections
func proxyTraffic(conn io.ReadWriter, link *Link, done chan bool) {
	defer close(link.ch)
	
	buffer := make([]byte, 32*1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			break
		}

		if _, err := link.Write(buffer[:n]); err != nil {
			break
		}
	}
	done <- true
}
