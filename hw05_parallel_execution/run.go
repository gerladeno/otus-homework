package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type safeInt struct {
	mx  sync.Mutex
	val int
}

func worker(jobs <-chan Task, errChan chan<- struct{}) {
	for job := range jobs {
		if err := job(); err != nil {
			errChan <- struct{}{}
		}
	}
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	var mainError error
	tasksChan := make(chan Task)

	var errChanLen int
	if m >= 0 {
		errChanLen = m
	} else {
		errChanLen = len(tasks)
	}
	errChan := make(chan struct{}, errChanLen)
	errCnt := safeInt{}
	go func() {
		defer close(tasksChan)
		for _, task := range tasks {
			select {
			case <-errChan:
				if m >= 0 {
					errCnt.mx.Lock()
					errCnt.val++
					if errCnt.val >= m {
						errCnt.mx.Unlock()
						mainError = ErrErrorsLimitExceeded
						return
					}
					errCnt.mx.Unlock()
				}
				tasksChan <- task
			default:
				tasksChan <- task
			}
		}
	}()
	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(tasksChan, errChan)
		}()
	}

	wg.Wait()
	return mainError
}
