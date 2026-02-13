package collector

import (
	dto "github.com/prometheus/client_model/go"
)

type MetricsCache struct {
    mu      sync.RWMutex
    // Map target_name -> metrics
    storage map[string][]*dto.MetricFamily
}

func NewMetricsCache() *MetricsCache {
    return &MetricsCache{
        storage: make(map[string][]*dto.MetricFamily),
    }
}

// Update updates the cache with the latest metrics for a specific target.
func (mc *MetricsCache) Update(target string, metrics []*dto.MetricFamily) {
    mc.mu.Lock()
    defer mc.mu.Unlock()
    mc.storage[target] = metrics
}

// Gather implements prometheus.Gatherer.
// It returns all stored metrics from all targets.
func (mc *MetricsCache) Gather() ([]*dto.MetricFamily, error) {
    mc.mu.RLock()
    defer mc.mu.RUnlock()

    var allMetrics []*dto.MetricFamily
    for _, m := range mc.storage {
        allMetrics = append(allMetrics, m...)
    }
    return allMetrics, nil
}