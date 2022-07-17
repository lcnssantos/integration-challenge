package concurrency

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type TaskResult struct {
	Result interface{}
	Err    error
	index  int
}

type Task = func() (interface{}, error)
type TaskResultMap map[int]TaskResult

type TaskInput struct {
	Task Task
	Tag  string
}

func ExecuteConcurrentTasks(tasks ...TaskInput) TaskResultMap {
	taskResult := make(TaskResultMap, 0)
	asyncTaskChannel := make(chan TaskResult)

	wg := sync.WaitGroup{}
	wg.Add(len(tasks))

	startTime := time.Now()

	for index, task := range tasks {
		go func(index int, task TaskInput, channel chan TaskResult) {
			startTime := time.Now()
			result, err := task.Task()

			log.
				Debug().
				Str("tag", task.Tag).
				Dur("duration", time.Since(startTime)).
				Msgf("task completed %s", task.Tag)

			if err != nil {
				log.Error().Err(err).Str("tag", task.Tag).Msg("task failed")
			}

			channel <- TaskResult{
				Result: result,
				Err:    err,
				index:  index,
			}
		}(index, task, asyncTaskChannel)
	}

	for i := 0; i < len(tasks); i++ {
		message := <-asyncTaskChannel
		taskResult[message.index] = message
		wg.Done()
	}

	wg.Wait()

	log.
		Debug().
		Dur("duration", time.Since(startTime)).
		Msg("all tasks completed")

	return taskResult
}
