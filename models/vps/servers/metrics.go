package servers

// Measure represents a single metric measurement.
type Measure struct {
	Granularity int     `json:"granularity"`
	Timestamp   int64   `json:"timestamp"`
	Value       float64 `json:"value"`
}

// MetricInfo contains metric data with measures.
type MetricInfo struct {
	Measures []Measure `json:"measures"`
	Name     string    `json:"name"`
}

// ServerMetricsRequest specifies query parameters for Metrics.
type ServerMetricsRequest struct {
	Type        string `json:"type,omitempty"`        // cpu, memory, disk, net, vgpu
	Granularity int    `json:"granularity,omitempty"` // seconds
	Start       int64  `json:"start,omitempty"`       // Unix timestamp
	Direction   string `json:"direction,omitempty"`   // incoming/outgoing for net
	RW          string `json:"rw,omitempty"`          // read/write for disk
}

// ServerMetricsResponse is the response from Metrics (array of MetricInfo).
type ServerMetricsResponse []MetricInfo
