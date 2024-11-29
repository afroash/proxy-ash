package metrics

type ConnectionStats struct {
	BytesSent     int64
	BytesReceived int64
	PacketsLost   int64
	LatencyMs     float64
	StartTime     int64
	EndTime       int64
}

type Stats struct {
	ActiveConnections int64
	TotalConnections  int64
	TotalBytes        int64
	AverageLatency    float64
}
