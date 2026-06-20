package pool

import (
	"context"
)

type Task interface {
	Execute(ctx context.Context) (interface{}, error)
}

type Job struct {
	ID   int
	Task Task
}

type Result struct {
	JobID  int
	Output interface{}
	Err    error
}
