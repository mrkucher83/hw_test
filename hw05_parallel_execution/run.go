package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrNegativeNumber      = errors.New("the number of errors or workers must be positive")
)

type Task func() error

type ErrCounter struct {
	mu    sync.Mutex
	value int
}

func (e *ErrCounter) Increase() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.value++
}

func (e *ErrCounter) Get() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.value
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m < 0 || n <= 0 {
		return ErrNegativeNumber
	}

	ch := make(chan Task)
	errCounter := ErrCounter{}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for _, task := range tasks {
			if errCounter.Get() == m {
				break
			}
			ch <- task
		}
		close(ch)
		wg.Done()
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			for task := range ch {
				if err := task(); err != nil {
					errCounter.Increase()
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if errCounter.Get() >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
