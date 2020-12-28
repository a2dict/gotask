package gotask

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestSubmit(t *testing.T) {
	w := NewWorker(2, 4)

	idx := int32(0)
	run := func() error {
		t := atomic.AddInt32(&idx, 1)
		time.Sleep(time.Duration(t) * time.Second)
		fmt.Printf("run: %v\n", t)
		return nil
	}

	for i := 0; i < 6; i++ {
		w.Submit(run)
		fmt.Printf("task submitted: %v\n", i)
	}

	w.Stop()
	w.Join()

	fmt.Println("exit")
}

func TestSubmitC(t *testing.T) {
	run := func() error {
		time.Sleep(10 * time.Second)
		return nil
	}

	timeout, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	w := NewWorker(1, 4)
	errC1 := w.SubmitC(timeout, run)
	errC2 := w.SubmitC(timeout, run)

	fmt.Println(<-errC1, <-errC2) // <nil>, context deadline exceeded

}
