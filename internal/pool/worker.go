package pool

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan Job, results chan<- Result) {
	for job := range jobs {
		fmt.Printf("[Worker %d] started processing job %d\n", id, job.ID)

		time.Sleep(500 * time.Millisecond)

		output := fmt.Sprintf("Processed payload: '%s'", job.Data)

		results <- Result{
			JobID:  job.ID,
			Output: output,
		}

		fmt.Printf("[Worker %d] finished job %d\n", id, job.ID)
	}
}
