package worker

import (
	"sync/atomic"
	"testing"
)

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool(3, 10)
	pool.Start()

	var counter int64
	totalTasks := 10

	for i := 0; i < totalTasks; i++ {
		pool.Submit(func() {
			atomic.AddInt64(&counter, 1)
		})
	}

	pool.Stop()

	if counter != int64(totalTasks) {
		t.Fatalf("excepted counter %d, got %d", totalTasks, counter)
	}
}
