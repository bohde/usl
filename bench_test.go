package usl_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/joshbohde/usl"
	"golang.org/x/sync/semaphore"
)

func ExampleBenchmark() {
	testing.Benchmark(func(b *testing.B) {
		client := http.Client{}

		usl.Bench(b, func(pb *usl.PB) {
			for pb.Next() {
				resp, err := client.Get("https://google.com")
				if err != nil {
					b.Fatal(err)
					return
				}
				_ = resp.Body.Close()
			}
		})

	})
}

func BenchmarkTest(b *testing.B) {
	s := semaphore.NewWeighted(10)

	usl.Bench(b, func(pb *usl.PB) {
		for pb.Next() {
			s.Acquire(context.Background(), 1)
			time.Sleep(30 * time.Millisecond)
			s.Release(1)
		}
	})
}
