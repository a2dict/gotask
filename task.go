package gotask

import (
	"context"
	"errors"
	"sync"
)

type runnable func() error

type task struct {
	ctx  context.Context
	errC chan error
	run  runnable
}

func (t *task) do() {
	select {
	case <-t.ctx.Done():
		t.errC <- t.ctx.Err()
	default:
		t.errC <- t.run()
	}
}

func NewWorker(poolSize, taskSize int) *worker {
	w := &worker{
		taskChan: make(chan *task, taskSize),
		stop:     make(chan struct{}),
	}

	for i := 0; i < poolSize; i++ {
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			for {
				select {
				case <-w.stop:
					return
				case t := <-w.taskChan:
					t.do()
				}
			}
		}()
	}

	return w
}

type worker struct {
	wg       sync.WaitGroup
	taskChan chan *task
	stop     chan struct{}
}

func (w *worker) SubmitC(ctx context.Context, run runnable) <-chan error {
	errC := make(chan error, 1)
	if ctx == nil {
		errC <- errors.New("nil context")
		return errC
	}

	t := &task{
		ctx:  ctx,
		errC: errC,
		run:  run,
	}
	w.taskChan <- t
	return errC
}

func (w *worker) Submit(run runnable) <-chan error {
	return w.SubmitC(context.Background(), run)
}

func (w *worker) Stop() {
	close(w.stop)
}

func (w *worker) Join() {
	w.wg.Wait()
}
