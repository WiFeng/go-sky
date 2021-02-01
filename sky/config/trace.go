package config

import "time"

// Trace ...
type Trace struct {
	Reporter Reporter
}

// Reporter ...
type Reporter struct {
	CollectorEndpoint   string
	LocalAgentHostPort  string
	BufferFlushInterval time.Duration
}
