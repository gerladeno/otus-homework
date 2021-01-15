package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(jobs <-chan Task, result chan<- struct{}) {
	for i := range jobs {
		if err := i(); err != nil {
			result <- struct{}{}
		}
	}
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, N int, M int) error {
	tasksChan := make(chan Task, len(tasks))
	errChan := make(chan struct{})
	go func() {
		tasksLoop: for _, task := range tasks {
			select {
			case tasksChan <- task:
			case <-errChan:
				break
			}

		}
	}()
	go func() {
		for i := 0; i < N; i++ {
			go worker(tasksChan, errChan)
		}
	}()
	return nil
}
