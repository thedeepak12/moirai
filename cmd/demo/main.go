package main

import (
	"fmt"
	"net/http"
	"time"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/thedeepak12/moirai/internal/pool"
)

type HTTPCheckTask struct {
	URL string
}
func (h HTTPCheckTask) Execute() (interface{}, error) {
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	start := time.Now()
	resp, err := client.Get(h.URL)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	duration := time.Since(start)
	
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("received bad status code: %d", resp.StatusCode)
	}
	return fmt.Sprintf("URL %s is healthy. Status: %d, Response Time: %s", h.URL, resp.StatusCode, duration), nil
}

type PasswordHashTask struct {
	Password string
	Salt     string
}
func (p PasswordHashTask) Execute() (interface{}, error) {
	if len(p.Password) < 8 {
		return nil, errors.New("Password is too short (min 8 chars)")
	}

	hash := sha256.New()
	data := []byte(p.Password + p.Salt)
	for i := 0; i < 10000; i++ {
		hash.Reset()
		hash.Write(data)
		data = hash.Sum(nil)
	}

	return hex.EncodeToString(data), nil
}

func main() {
	numWorkers := 3

	tasks := []pool.Task{
		HTTPCheckTask{URL: "https://www.google.com"},
		HTTPCheckTask{URL: "https://go.dev"},

		PasswordHashTask{Password: "v5x!wKZmp2VF9HKYyeRd$^VU", Salt: "ftsz3pLHiS0Lid/17zkB4wddmsDPsg=="},
		PasswordHashTask{Password: "123456", Salt: "oKuNU/4X5bdWsEu4OnZlfoPgEPIk8g=="},
	}

	p := pool.NewPool(numWorkers, len(tasks))
	p.Start()

	go func() {
		for i, t := range tasks {
			p.Submit(pool.Job{
				ID:   i + 1,
				Task: t,
			})
		}

		p.Wait()
	}()

	fmt.Println("Starting to receive results...")
	for result := range p.Results() {
		if result.Err != nil {
			fmt.Printf("[Main] Job %d failed with error: %v\n", result.JobID, result.Err)
		} else {
			fmt.Printf("[Main] Job %d succeeded: %v\n", result.JobID, result.Output)
		}
	}

	fmt.Println("All done! Exiting.")
}
