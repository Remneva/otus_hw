package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	workCh := make(chan Task)

	wg := sync.WaitGroup{}
	wg.Add(n)
	var z int32
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range workCh {
				err := task()
				if err != nil {
					fmt.Println(atomic.AddInt32(&z, 1))
				}
			}
		}()
	}
	for _, task := range tasks {
		if atomic.LoadInt32(&z) >= int32(m) {
			break
		}
		workCh <- task
	}

	close(workCh)
	wg.Wait()
	if z >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
