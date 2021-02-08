package viewer

import "github.com/pkg/errors"

// Metric represents a metric value tied to a Viewer
type Metric struct {
	Name   string
	Amount int
}

// MetricStore is a key/value store for viewer.Metric
type MetricStore interface {
	Fetch(key string) (Metric, error)
	Put(key string, value Metric) error
}

// InMemoryMetricStore implements MetricStore in memory
type InMemoryMetricStore struct {
	store map[string]Metric
}

func NewInMemoryStore() *InMemoryMetricStore {
	return &InMemoryMetricStore{
		store: make(map[string]Metric),
	}
}

// Fetch fetches the Metric at the given key, or returns an error
// if key not found
func (m *InMemoryMetricStore) Fetch(key string) (Metric, error) {
	metric, ok := m.store[key]

	if !ok {
		return Metric{}, errors.Errorf("Key not found: %s", key)
	}
	return metric, nil
}

// Put adds a new key/value pair to the store
func (m *InMemoryMetricStore) Put(key string, value Metric) error {
	m.store[key] = value
	return nil
}
