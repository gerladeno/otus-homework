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

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, N int, M int) error {
	var mainError error
	tasksChan := make(chan Task)
	errChan := make(chan struct{}, M)
	errCnt := safeInt{}
	go func() {
	tasksLoop:
		for _, task := range tasks {
			select {
			case <-errChan:
				errCnt.mx.Lock()
				errCnt.val++
				if errCnt.val >= M {
					errCnt.mx.Unlock()
					mainError = ErrErrorsLimitExceeded
					break tasksLoop
				}
				errCnt.mx.Unlock()
				tasksChan <- task
			default:
				tasksChan <- task
			}
		}
		close(tasksChan)
	}()
	wg := sync.WaitGroup{}
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(tasksChan, errChan)
		}()
	}

	wg.Wait()
	return mainError
}
