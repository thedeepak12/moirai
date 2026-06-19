package pool

type Job struct {
	ID   int
	Data string
}

type Result struct {
	JobID  int
	Output string
}
