package pool

import (
	"sync/atomic"
)

type Metrics struct {
	submitted int64
	completed int64
	failed    int64
}

type Progress struct {
	Submitted int64
	Completed int64
	Failed    int64
}

func (m *Metrics) IncSubmitted() {
	atomic.AddInt64(&m.submitted, 1)
}

func (m *Metrics) IncCompleted() {
	atomic.AddInt64(&m.completed, 1)
}
func (m *Metrics) IncFailed() {
	atomic.AddInt64(&m.failed, 1)
}

func (m *Metrics) Snapshot() Progress {
	return Progress{
		Submitted: atomic.LoadInt64(&m.submitted),
		Completed: atomic.LoadInt64(&m.completed),
		Failed:    atomic.LoadInt64(&m.failed),
	}
}
