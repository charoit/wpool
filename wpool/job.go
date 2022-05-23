package wpool

import (
	"context"
)

type ExecFn func(ctx context.Context, args interface{}) (interface{}, error)

type Result struct {
	Value    interface{}
	Err      error
	Metadata map[string]interface{}
}

type Job struct {
	Args     interface{}
	ExecFn   ExecFn
	Metadata map[string]interface{}
}

func (j Job) execute(ctx context.Context) Result {
	value, err := j.ExecFn(ctx, j.Args)
	if err != nil {
		return Result{
			Err:      err,
			Metadata: j.Metadata,
		}
	}

	return Result{
		Value:    value,
		Metadata: j.Metadata,
	}
}
