package executor

import (
	"errors"
	"fmt"
	"reflect"
)

type Task struct {
	handler any
	args    []reflect.Value
}

func NewTask(handler any, inputArgs ...any) (*Task, error) {
	nArgs := len(inputArgs)
	parsedHandler, err := validateFunc(handler, nArgs)
	if err != nil {
		return nil, err
	}
	var args = make([]reflect.Value, 0, nArgs)
	for _, inputArg := range inputArgs {
		args = append(args, reflect.ValueOf(inputArg))
	}
	return &Task{handler: parsedHandler, args: args}, nil
}

func validateFunc(handler any, nArgs int) (any, error) {
	f := reflect.Indirect(reflect.ValueOf(handler))
	if f.Kind() != reflect.Func {
		return f, fmt.Errorf("%T must be a function", f)
	}
	numIn := reflect.ValueOf(handler).Type().NumIn()
	if nArgs < numIn {
		return nil, errors.New("Call with too few args")
	} else if nArgs > numIn {
		return nil, errors.New("Call with too many args")
	}
	return f, nil
}
