package executor

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

type Executor struct {
	waitGroup *sync.WaitGroup
	tasks     chan *Task
	workers   int
}

type Config struct {
	QueueSize  int
	NumWorkers int
}

func New(config Config) (*Executor, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}
	numWorkers := config.NumWorkers
	if numWorkers == 0 {
		numWorkers = runtime.NumCPU()
	}
	queueSize := config.QueueSize
	if queueSize == 0 {
		queueSize = 2 * numWorkers
	}
	executor := &Executor{
		waitGroup: new(sync.WaitGroup),
		tasks:     make(chan *Task, queueSize),
		workers:   numWorkers,
	}
	executor.initWorker()
	return executor, nil
}

func (config *Config) validate() error {
	if config.QueueSize < 0 {
		return fmt.Errorf("%T must positive", "QueueSize")
	}
	if config.NumWorkers < 0 {
		return fmt.Errorf("%T must greater than 0", "NumWorkers")
	}
	return nil
}

func (executor *Executor) initWorker() {
	for i := 0; i < executor.workers; i++ {
		go executor.runWorker()
	}
}

func (executor *Executor) runWorker() {
	for {
		task, ok := <-executor.tasks
		if !ok {
			break
		}

		fn := task.handler.(reflect.Value)
		_ = fn.Call(task.args)

		executor.waitGroup.Done()
	}
}

func (executor *Executor) Wait() {
	executor.waitGroup.Wait()
}

func (executor *Executor) Close() {
	executor.Wait()
	close(executor.tasks)
}

func (executor *Executor) Publish(handler interface{}, args ...interface{}) error {
	task, err := NewTask(handler, args...)
	if err != nil {
		return err
	}
	executor.PublishTask(task)
	return nil
}
func (executor *Executor) PublishTask(task *Task) {
	executor.waitGroup.Add(1)
	executor.tasks <- task
}
