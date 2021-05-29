package viewer

import (
	"bytes"
	"encoding/json"
	"strings"
)

const (
	// LastSub is the Cache key for the last subscriber
	LastSub = "lastSub"

	// LastGiftSub is the Cache key for the last user to gift subscriptions
	LastGiftSub = "lastGiftSub"

	// LastBits is the Cache key for the last bits donation
	LastBits = "lastBits"
)

// Metric represents a metric value tied to a Viewer
type Metric struct {
	Name      string `json:"name"`
	Recipient string `json:"recipient,omitempty"`
	Amount    int    `json:"amount"`
}

// String representation of a Metric
func (m Metric) String() string {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(m)
	return strings.TrimSuffix(buf.String(), "\n")
}

// FromString attempts to convert a string (created by Metric.String()) to a Metric
func FromString(line string) (Metric, error) {
	var metric Metric

	if len(line) <= 0 {
		return Metric{}, nil
	}

	err := json.Unmarshal([]byte(line), &metric)
	if err != nil {
		return Metric{}, err
	}

	return metric, nil
}
