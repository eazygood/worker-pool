package internal

import (
	"fmt"
	"time"
)

type HttpResponseMetrics struct {
	Label    string
	Domain   string
	Duration time.Duration
	SizeInKb float64
}

func (h HttpResponseMetrics) String() string {
	return fmt.Sprintf("Method: %v\nDomail: %v\nDuration: %v\nSize: %vKb\n ------", h.Label, h.Domain, h.Duration, h.SizeInKb)
}
