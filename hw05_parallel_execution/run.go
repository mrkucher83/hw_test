package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrNegativeNumber      = errors.New("the number of errors or workers must be positive")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m < 0 || n <= 0 {
		return ErrNegativeNumber
	}

	ch := make(chan Task)
	var ErrCounter int32

	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			for task := range ch {
				if err := task(); err != nil {
					atomic.AddInt32(&ErrCounter, 1)
				}
			}
			wg.Done()
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&ErrCounter) >= int32(m) {
			break
		}
		ch <- task
	}
	close(ch)
	wg.Wait()

	if ErrCounter >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
