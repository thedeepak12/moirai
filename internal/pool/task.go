package pool

type Task interface {
	Execute() (interface{}, error)
}

type Job struct {
	ID   int
	Task Task
}

type Result struct {
	JobID  int
	Output interface{}
	Err error
}
