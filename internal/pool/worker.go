package pool

import (
	"fmt"	
)

func worker(id int, jobs <-chan Job, results chan<- Result) {
	for job := range jobs {
		fmt.Printf("[Worker %d] started processing job %d\n", id, job.ID)

		output, err := job.Task.Execute()

		results <- Result{
			JobID:  job.ID,
			Output: output,
			Err: err,
		}

		if err != nil {
			fmt.Printf("[Worker %d] Job %d failed: %v\n", id, job.ID, err)
		} else {
			fmt.Printf("[Worker %d] finished job %d successfully\n", id, job.ID)
		}
	}
}
