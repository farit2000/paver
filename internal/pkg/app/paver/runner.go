package paver

import (
	"errors"
	"sync"
	"sync/atomic"
)

type Task struct {
	ID       string
	TaskFunc func() (interface{}, error)
}
type TaskResult struct {
	ID    string
	Value interface{}
	Error error
}

type Worker struct {
	id      int
	task    chan func()
	stop    chan bool
	stopped chan bool
}

func NewWorker(id int, pool chan func()) *Worker {
	worker := &Worker{
		id:      id,
		task:    pool,
		stop:    make(chan bool),
		stopped: make(chan bool),
	}

	go func() {
		for {
			select {
			case task := <-worker.task:
				task()
			case <-worker.stop:
				close(worker.stopped)
				return
			}
		}
	}()

	return worker
}

func (w *Worker) Stop() {
	close(w.stop)
	<-w.stopped
}

type WorkersPool struct {
	workerQueue    chan func()
	workers        []*Worker
	wg             sync.WaitGroup
	shutdownInited atomic.Bool
}

func NewWorkersPool(numWorkers int) *WorkersPool {
	pool := &WorkersPool{
		workerQueue: make(chan func(), numWorkers),
		workers:     make([]*Worker, numWorkers),
	}

	for i := 0; i < numWorkers; i++ {
		pool.workers[i] = NewWorker(i, pool.workerQueue)
	}

	return pool
}

func (p *WorkersPool) Submit(task Task, resultChan chan<- TaskResult) error {
	if p.shutdownInited.Load() {
		return errors.New("try to submit task to stopped pool")
	}
	p.wg.Add(1)
	taskWrapper := func() {
		defer p.wg.Done()
		result, err := task.TaskFunc()
		resultChan <- TaskResult{
			ID:    task.ID,
			Value: result,
			Error: err,
		}
	}
	p.workerQueue <- taskWrapper

	return nil
}

func (p *WorkersPool) Shutdown() {
	p.shutdownInited.Store(true)
	p.wg.Wait()
	for _, worker := range p.workers {
		worker.Stop()
	}
	close(p.workerQueue)
}
