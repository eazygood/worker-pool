package internal

import "context"

type ExecFunction func(ctx context.Context, args interface{}) (interface{}, error)

type JobDescriptor struct {
	JobType  string
	Metadata map[string]interface{}
}

type Job struct {
	ID       int
	ExecFunc ExecFunction
	Args     interface{}
}

type Result struct {
	Value interface{}
	Err   error
	JobID int
}

func (j Job) execute(ctx context.Context) Result {
	value, err := j.ExecFunc(ctx, j.Args)
	if err != nil {
		return Result{
			Err:   err,
			JobID: j.ID,
		}
	}

	return Result{
		Value: value,
		JobID: j.ID,
	}
}
