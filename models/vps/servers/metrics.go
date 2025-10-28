package servers

// ServerMetricsRequest specifies query parameters for Metrics.
type ServerMetricsRequest struct {
	Type        string // cpu, memory, disk, network
	Start       int64  // Unix timestamp
	End         int64  // Unix timestamp
	Granularity int    // seconds
}

// ServerMetricsResponse is the response from Metrics.
type ServerMetricsResponse struct {
	Type   string        `json:"type"`
	Series []MetricPoint `json:"series"`
}

// MetricPoint is a single time-series data point.
type MetricPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}
