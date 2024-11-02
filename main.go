package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp4", "0.0.0.0:42069")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Proxy up and listening on 42069")
	defer listener.Close()

	for {
		downstreamConn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(downstreamConn)
	}
}

type link struct {
	ch      chan []byte
	latency time.Duration
}

func (l *link) Read(p []byte) (n int, err error) {
	data, ok := <-l.ch
	if !ok {
		return 0, io.EOF
	}
	time.Sleep(l.latency)
	n = copy(p, data)
	return n, nil
}

func (l *link) Write(p []byte) (n int, err error) {
	data := make([]byte, len(p))
	copy(data, p)
	l.ch <- data
	return len(p), nil
}

func handleConnection(downstreamConn net.Conn) {
	defer downstreamConn.Close()

	// Connect to upstream server
	upstreamConn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Failed to connect to upstream:", err)
		return
	}
	defer upstreamConn.Close()

	// Create links with latency
	upstreamLink := &link{
		ch:      make(chan []byte, 1024),
		latency: time.Second * 2,
	}
	downstreamLink := &link{
		ch:      make(chan []byte, 1024),
		latency: time.Second * 2,
	}

	// Handle bidirectional communication
	done := make(chan bool, 2)

	// Downstream -> Upstream (through upstreamLink)
	go func() {
		defer close(upstreamLink.ch)
		buf := make([]byte, 32*1024)
		for {
			n, err := downstreamConn.Read(buf)
			if err != nil {
				break
			}
			// Write to link
			if _, err := upstreamLink.Write(buf[:n]); err != nil {
				break
			}
		}
		done <- true
	}()

	// Upstream -> Downstream (through downstreamLink)
	go func() {
		defer close(downstreamLink.ch)
		buf := make([]byte, 32*1024)
		for {
			n, err := upstreamConn.Read(buf)
			if err != nil {
				break
			}
			// Write to link
			if _, err := downstreamLink.Write(buf[:n]); err != nil {
				break
			}
		}
		done <- true
	}()

	// Read from links and forward
	go func() {
		buf := make([]byte, 32*1024)
		for {
			n, err := upstreamLink.Read(buf)
			if err != nil {
				break
			}
			if _, err := upstreamConn.Write(buf[:n]); err != nil {
				break
			}
		}
		done <- true
	}()

	go func() {
		buf := make([]byte, 32*1024)
		for {
			n, err := downstreamLink.Read(buf)
			if err != nil {
				break
			}
			if _, err := downstreamConn.Write(buf[:n]); err != nil {
				break
			}
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	for i := 0; i < 4; i++ {
		<-done
	}
}
