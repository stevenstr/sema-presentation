package sema3_test

import (
	"testing"
	"sync"
	"golang.org/x/sync/semaphore"
	"context"
	"time"
	"sync/atomic"
)

type request struct{}

var count int64

func handle(r request) {
	atomic.AddInt64(&count, 1)
	time.Sleep(time.Microsecond)
}

func process(maxHandlers int, reqs []request) {
	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(int64(maxHandlers))

	for _, r := range reqs {
		wg.Add(1)
		go func(r request) {
			sem.Acquire(context.Background(), 1)
			handle(r)
			sem.Release(1)
			wg.Done()
		}(r)
	}
	wg.Wait()
}

func benchmarkServe(b *testing.B, n int) {
	count = 0
	reqs := make([]request, b.N)
	process(n, reqs)
	if int(count) != b.N {
		b.Errorf("number of messages handled doesn't match, wanted: '%d' but received: '%d'", b.N, count)

	}
}

func BenchmarkServe10(b *testing.B) {
	benchmarkServe(b, 10)
}

func BenchmarkServe100(b *testing.B) {
	benchmarkServe(b, 100)
}

func BenchmarkServe1000(b *testing.B) {
	benchmarkServe(b, 1000)
}

func BenchmarkServe10000(b *testing.B) {
	benchmarkServe(b, 10000)
}

func BenchmarkServe100000(b *testing.B) {
	benchmarkServe(b, 100000)
}
