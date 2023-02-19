package executor

import (
	assrt "github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestValidateDefaultConfig(t *testing.T) {
	assert := assrt.New(t)

	executor, err := New(Config{})
	defer executor.Close()

	assert.Equal(runtime.NumCPU(), executor.workers)
	assert.NotNil(executor.tasks)
	assert.Equal(runtime.NumCPU()*2, cap(executor.tasks))
	assert.NotNil(executor.waitGroup)
	assert.Nil(err)
}

func TestValidateConfig(t *testing.T) {
	assert := assrt.New(t)

	executor, err := New(Config{QueueSize: 20, NumWorkers: 2})
	defer executor.Close()

	assert.Equal(2, executor.workers)
	assert.NotNil(executor.tasks)
	assert.Equal(20, cap(executor.tasks))
	assert.NotNil(executor.waitGroup)
	assert.Nil(err)
}

func TestPublishJobSuccess(t *testing.T) {
	assert := assrt.New(t)

	executor, err := New(Config{})
	assert.Nil(err)

	value := 1
	err = executor.Publish(
		func(input int) {
			assert.Equal(value, input)
		},
		value)
	assert.Nil(err)
	executor.Wait()
}

func TestPublishJobFail(t *testing.T) {
	assert := assrt.New(t)

	executor, err := New(Config{})
	assert.Nil(err)

	err = executor.Publish(
		func(input int) {
			assert.Equal(1, input)
		},
		1, 1)

	assert.NotNil(err)
	assert.Equal(err.Error(), "Call with too many args")
	executor.Wait()
}
