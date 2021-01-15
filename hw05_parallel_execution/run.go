package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type safeCnt struct {
	cnt int
	mx  sync.Mutex
}

func worker(jobs <-chan Task, result chan<- struct{}) {
	for i := range jobs {
		if err := i(); err != nil {
			result <- struct{}{}
		}
	}
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, N int, M int) error {
	wg := sync.WaitGroup{}
	errs := safeCnt{}
	tasksChan := make(chan Task)
	queue := make(chan struct{}, N)
	errChan := make(chan struct{})
	allDone := make(chan struct{})
	go func() {
		for _, task := range tasks {
			tasksChan<- task
			select {
			case <-errChan:
				break
			}
		}
	}()
	for i := 0; i < N; i++ {
		queue <- struct{}{}
		go func() {
			wg.Add(1)
			defer wg.Done()
			if err := task(); err != nil {
				errs.mx.Lock()
				errs.cnt++
				if errs.cnt >= M {
					errChan <- struct{}{}
					return
				}
				errs.mx.Unlock()
			}
			<-queue
		}()
	}

	go func() {
		wg.Wait()
		allDone <- struct{}{}
	}()
	select {
	case <-errChan:
		return ErrErrorsLimitExceeded
	case <-allDone:
		return nil
	}
}
