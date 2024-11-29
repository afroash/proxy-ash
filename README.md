# Network Proxy with Traffic Simulation

A configurable TCP proxy server that can simulate various network conditions like latency, packet loss, and bandwidth limitations. Perfect for testing application behavior under different network scenarios.

## Features

- **TCP Proxy**: Forward TCP traffic between client and upstream server
- **Network Condition Simulation**:
  - Configurable latency (min/max delay)
  - Packet loss simulation
  - Bandwidth throttling
- **Metrics Collection**: Monitor proxy performance and connection statistics
- **YAML Configuration**: Easy configuration of all proxy settings

## Installation

```bash
# Clone the repository
git clone https://github.com/afroash/proxy-ash.git
cd proxy-ash

# Build the proxy
go build -o proxy ./cmd/proxy

# Build the test server (optional)
go build -o testserver ./cmd/testserver
```

## Configuration

Create a `config.yaml` file with your desired settings:

```yaml
ListenAddr: ":8080"
UpstreamAddr: "localhost:9090"

# Network condition simulation settings
latency:
  enabled: true
  min_ms: 50
  max_ms: 150
  duration: "100ms"

packet_loss:
  enabled: true
  percentage: 5.0

bandwidth:
  enabled: true
  limit_kbps: 1024  # 1 Mbps

# Metrics collection settings
metrics:
  enabled: true
  path: "/metrics"
```

## Usage

1. Start the proxy:
```bash
./proxy -config config.yaml
```

2. (Optional) Start the test server for testing:
```bash
./testserver
```

3. The proxy will now:
   - Listen on the configured port (default: 8080)
   - Forward traffic to the upstream server
   - Apply configured network conditions
   - Collect metrics (if enabled)

## Testing

You can test the proxy using curl or any HTTP client:

```bash
# Test through proxy
curl -v http://localhost:8080/test

# Direct test to upstream
curl -v http://localhost:9090/test
```

Compare the response times and behavior between direct and proxied requests to observe the simulated network conditions.

## Metrics

When metrics are enabled, you can access them at the configured metrics path:

```bash
curl http://localhost:8080/metrics
```

Available metrics include:
- Active connections
- Total bytes transferred
- Average latency
- Packet loss statistics

## Project Structure

```
.
├── api/
│   └── server.go         # Main proxy server implementation
├── cmd/
│   ├── proxy/           # Proxy server entry point
│   └── testserver/      # Test server for demonstration
├── internal/
│   ├── config/          # Configuration handling
│   ├── metrics/         # Metrics collection
│   ├── proxy/           # Proxy core functionality
│   └── simulator/       # Network condition simulation
├── config.yaml          # Configuration file
└── README.md           # This file
```

## Contributing

Feel free to submit issues, fork the repository, and create pull requests for any improvements.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
