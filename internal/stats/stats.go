package stats

import (
	"sort"
	"time"
)

type Summary struct {
	Count int

	Min time.Duration
	Max time.Duration

	P50 time.Duration
	P95 time.Duration
	P99 time.Duration

	Avg time.Duration
}

func Calculate(latencies []time.Duration) Summary {
	n := len(latencies)
	if n == 0 {
		return Summary{}
	}

	sorted := make([]time.Duration, n)
	copy(sorted, latencies)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	sum := time.Duration(0)
	for _, v := range sorted {
		sum += v
	}

	return Summary{
		Count: n,
		Min:   sorted[0],
		Max:   sorted[n-1],
		P50:   percentile(sorted, 50),
		P95:   percentile(sorted, 95),
		P99:   percentile(sorted, 99),
		Avg:   sum / time.Duration(n),
	}
}

func percentile(sorted []time.Duration, p int) time.Duration {
	if len(sorted) == 0 {
		return 0
	}

	idx := int(float64(len(sorted)-1) * float64(p) / 100.0)
	return sorted[idx]
}
