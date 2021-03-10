package viewer

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
)

// Metric represents a metric value tied to a Viewer
type Metric struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

// String representation of a Metric
func (m Metric) String() string {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(m)
	return buf.String()
}

// FromCache attempts to convert a string to a Metric
func FromCache(line string) (Metric, error) {
	var metric Metric

	if len(line) <= 0 {
		return Metric{}, errors.New("Empty line")
	}

	err := json.Unmarshal([]byte(line), &metric)
	if err != nil {
		return Metric{}, err
	}

	return metric, nil
}
