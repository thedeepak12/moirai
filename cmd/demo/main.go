package main

import (
	"fmt"

	"github.com/thedeepak12/moirai/internal/pool"
)

func main() {
	numWorkers := 3
	numJobs := 10

	p := pool.NewPool(numWorkers, numJobs)
	p.Start()

	go func() {
		for i := 1; i <= numJobs; i++ {
			p.Submit(pool.Job{
				ID:   i,
				Data: fmt.Sprintf("Data packet %d", i),
			})
		}

		p.Wait()
	}()

	fmt.Println("Starting to receive results...")
	for result := range p.Results() {
		fmt.Printf("[Main] Received Job %d output: %s\n", result.JobID, result.Output)
	}

	fmt.Println("All done! Exiting.")
}
