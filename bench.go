package usl

import (
	"sync"
	"testing"
	"time"
)

type PB struct {
	mu     sync.Mutex
	counts int64
	n      int64
}

func (pb *PB) Next() bool {
	pb.mu.Lock()
	if pb.counts < pb.n {
		pb.counts++
		pb.mu.Unlock()
		return true
	}

	pb.mu.Unlock()
	return false
}

func Bench(b *testing.B, f func(pb *PB)) (Model, error) {
	threadCounts := []int{1, 2, 4, 8, 16, 32, 64}
	points := make([][2]float64, len(threadCounts))

	n := (b.N / len(threadCounts))
	if n < 1 {
		n = 1
	}
	b.N = n * len(threadCounts)

	for i, c := range threadCounts {
		points[i] = [2]float64{float64(c), 0}
		throughput := &points[i][1]

		start := time.Now()

		pb := &PB{
			n: int64(n),
		}

		b.StartTimer()

		wg := sync.WaitGroup{}
		wg.Add(c)
		for i := 0; i < c; i++ {
			go func() {
				f(pb)
				wg.Done()
			}()
		}

		wg.Wait()
		b.StopTimer()

		seconds := float64(time.Since(start)) / float64(time.Second)
		*throughput = float64(pb.counts) / seconds

	}

	m, err := Fit(points)

	b.ReportMetric(m.Sigma, "sigma")
	b.ReportMetric(m.Kappa, "kappa")
	b.ReportMetric(m.Lambda, "lambda")
	b.ReportMetric(m.R2, "r2")
	b.ReportMetric(m.MaxConcurrency(), "max-concurrency")
	b.ReportMetric(m.MaxThroughput(), "max-throughput")

	return m, err
}
