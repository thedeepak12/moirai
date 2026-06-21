package pool

import (
	"context"
	"fmt"
)

func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result, metrics *Metrics, progress chan<- Progress, workerStop <-chan struct{}) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[Worker %d] Context cancelled, shutting down.\n", id)
			return

		case <-workerStop:
			fmt.Printf("[Worker %d] Scaled down, shutting down.\n", id)
			return

		case job, ok := <-jobs:
			if !ok {
				return
			}

			if err := ctx.Err(); err != nil {
				results <- Result{JobID: job.ID, Err: err}
				return
			}

			fmt.Printf("[Worker %d] started processing job %d\n", id, job.ID)

			output, err := job.Task.Execute(ctx)
			if err != nil {
				metrics.IncFailed()
			} else {
				metrics.IncCompleted()
			}

			select {
			case progress <- metrics.Snapshot():
			default:
			}

			results <- Result{
				JobID:  job.ID,
				Output: output,
				Err:    err,
			}

			if err != nil {
				fmt.Printf("[Worker %d] Job %d failed: %v\n", id, job.ID, err)
			} else {
				fmt.Printf("[Worker %d] finished job %d successfully\n", id, job.ID)
			}
		}
	}
}
