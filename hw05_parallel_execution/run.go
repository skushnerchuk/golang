package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded   = errors.New("errors limit exceeded")
	ErrIncorrectWorkersCount = errors.New("workers count must be > 0")
)

type Task func() error

func Run(tasks []Task, maxWorkersCount, maxErrorsCount int) error {
	if maxWorkersCount <= 0 {
		return ErrIncorrectWorkersCount
	}
	wg := sync.WaitGroup{}
	var errCounter int32
	maxErrorsCount32 := int32(maxErrorsCount)

	taskChannel := make(chan Task)

	wg.Add(maxWorkersCount)
	for i := 0; i < maxWorkersCount; i++ {
		go func() {
			defer wg.Done()
			for task := range taskChannel {
				if err := task(); err != nil {
					atomic.AddInt32(&errCounter, 1)
				}
			}
		}()
	}
	for _, t := range tasks {
		if atomic.LoadInt32(&errCounter) >= maxErrorsCount32 && maxErrorsCount32 > 0 {
			break
		}
		if t != nil {
			taskChannel <- t
		}
	}
	close(taskChannel)
	wg.Wait()
	if maxErrorsCount32 > 0 && (atomic.LoadInt32(&errCounter) >= maxErrorsCount32) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
