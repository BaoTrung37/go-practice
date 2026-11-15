package worker

import "sync"

type Task func()

type WorkerPool struct {
	numWorkers int
	taskQueue  chan Task
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int, queueSize int) *WorkerPool {
	return &WorkerPool{numWorkers: numWorkers, taskQueue: make(chan Task, queueSize)}
}

func (p *WorkerPool) Start() {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			for task := range p.taskQueue {
				task()
			}
		}(i)
	}
}

func (p *WorkerPool) Submit(t Task) {
	p.taskQueue <- t
}

func (p *WorkerPool) Stop() {
	close(p.taskQueue)
	p.wg.Wait()
}
