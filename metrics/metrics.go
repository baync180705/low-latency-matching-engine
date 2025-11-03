package metrics

import (
	"sort"
	"sync"
	"time"

	types "github.com/baync180705/low-latency-matching-engine/types"
)

type Snapshot struct {
	OrdersReceived         int64   `json:"orders_received"`
	OrdersMatched          int64   `json:"orders_matched"`
	OrdersCancelled        int64   `json:"orders_cancelled"`
	OrdersInBook           int64   `json:"orders_in_book"`
	TradesExecuted         int64   `json:"trades_executed"`
	LatencyP50Ms           float64 `json:"latency_p50_ms"`
	LatencyP99Ms           float64 `json:"latency_p99_ms"`
	LatencyP999Ms          float64 `json:"latency_p999_ms"`
	ThroughputOrdersPerSec float64 `json:"throughput_orders_per_sec"`
}

type metricsState struct {
	mu sync.RWMutex

	ordersReceived  int64
	ordersMatched   int64
	ordersCancelled int64
	tradesExecuted  int64

	latencies []float64 // in ms

	startTime time.Time
}

var state = &metricsState{
	latencies: []float64{},
	startTime: time.Now(),
}

// I initialize start time when the package is loaded.
func Init() {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.startTime = time.Now()
	state.latencies = state.latencies[:0]
	state.ordersReceived = 0
	state.ordersMatched = 0
	state.ordersCancelled = 0
	state.tradesExecuted = 0
}

// IncOrdersReceived increments the received-orders counter.
func IncOrdersReceived() {
	state.mu.Lock()
	state.ordersReceived++
	state.mu.Unlock()
}

// IncOrdersCancelled increments cancelled-orders counter.
func IncOrdersCancelled() {
	state.mu.Lock()
	state.ordersCancelled++
	state.mu.Unlock()
}

// AddTradesExecuted increases the total number of trade records created.
func AddTradesExecuted(n int) {
	if n <= 0 {
		return
	}
	state.mu.Lock()
	state.tradesExecuted += int64(n)
	state.mu.Unlock()
}

// AddOrdersMatched increases the matched-orders counter (unique orders touched by trades).
func AddOrdersMatched(n int) {
	if n <= 0 {
		return
	}
	state.mu.Lock()
	state.ordersMatched += int64(n)
	state.mu.Unlock()
}

// AddLatency records a single pipeline latency in milliseconds.
func AddLatency(ms float64) {
	state.mu.Lock()
	state.latencies = append(state.latencies, ms)
	state.mu.Unlock()
}

// percentile computes the p-th percentile (p in [0,100]) on a copy of the lat slice.
// If lat is empty, returns 0.
func percentile(lat []float64, p float64) float64 {
	n := len(lat)
	if n == 0 {
		return 0
	}
	sort.Float64s(lat)
	if p <= 0 {
		return lat[0]
	}
	if p >= 100 {
		return lat[n-1]
	}
	// Linear interpolation between nearest ranks
	pos := (p / 100.0) * float64(n-1)
	lower := int(pos)
	upper := lower + 1
	if upper >= n {
		return lat[lower]
	}
	frac := pos - float64(lower)
	return lat[lower]*(1-frac) + lat[upper]*frac
}

// GetSnapshot builds a Snapshot using internal counters and a live registry scan.
func GetSnapshot(registry *types.Regsitry) Snapshot {
	// Read counters
	state.mu.RLock()
	ordersReceived := state.ordersReceived
	ordersMatched := state.ordersMatched
	ordersCancelled := state.ordersCancelled
	tradesExecuted := state.tradesExecuted
	latCopy := make([]float64, len(state.latencies))
	copy(latCopy, state.latencies)
	startTime := state.startTime
	state.mu.RUnlock()

	// Compute orders_in_book by scanning registry (count of active orders)
	var ordersInBook int64 = 0
	if registry != nil {
		registry.Mu.RLock()
		for _, book := range registry.Books {
			book.Mu.RLock()
			for _, ord := range book.OrderIDMap {
				if !ord.IsComplete && !ord.IsCancelled {
					ordersInBook++
				}
			}
			book.Mu.RUnlock()
		}
		registry.Mu.RUnlock()
	}

	// Compute latency percentiles
	p50 := percentile(latCopy, 50.0)
	p99 := percentile(latCopy, 99.0)
	p999 := percentile(latCopy, 99.9)

	// Throughput = total orders received / uptimeSeconds
	uptime := time.Since(startTime).Seconds()
	var throughput float64
	if uptime <= 0 {
		throughput = 0
	} else {
		throughput = float64(ordersReceived) / uptime
	}

	return Snapshot{
		OrdersReceived:         ordersReceived,
		OrdersMatched:          ordersMatched,
		OrdersCancelled:        ordersCancelled,
		OrdersInBook:           ordersInBook,
		TradesExecuted:         tradesExecuted,
		LatencyP50Ms:           p50,
		LatencyP99Ms:           p99,
		LatencyP999Ms:          p999,
		ThroughputOrdersPerSec: throughput,
	}
}
