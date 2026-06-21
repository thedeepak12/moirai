package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/thedeepak12/moirai/internal/pool"
)

type HTTPCheckTask struct {
	URL string
}

func (h HTTPCheckTask) Execute(ctx context.Context) (interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", h.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := http.Client{}
	start := time.Now()
	resp, err := client.Do(req)
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

func (p PasswordHashTask) Execute(ctx context.Context) (interface{}, error) {
	if len(p.Password) < 8 {
		return nil, errors.New("Password is too short (min 8 chars)")
	}

	hash := sha256.New()
	data := []byte(p.Password + p.Salt)
	for i := 0; i < 10000; i++ {
		if i%1000 == 0 {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
		}

		hash.Reset()
		hash.Write(data)
		data = hash.Sum(nil)
	}

	return hex.EncodeToString(data), nil
}

func drawProgressBar(completed, total int64) string {
	width := 20

	if total == 0 {
		return "[" + strings.Repeat(" ", width) + "] 0%"
	}

	percent := float64(completed) / float64(total)
	filledLength := int(percent * float64(width))

	bar := strings.Repeat("█", filledLength) + strings.Repeat(" ", width-filledLength)
	return fmt.Sprintf("[%s] %.0f%%", bar, percent*100)
}

func main() {
	initialWorkers := 2

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tasks := []pool.Task{
		HTTPCheckTask{URL: "https://www.google.com"},
		HTTPCheckTask{URL: "https://go.dev"},
		HTTPCheckTask{URL: "https://github.com"},

		PasswordHashTask{Password: "v5x!wKZmp2VF9HKYyeRd$^VU", Salt: "ftsz3pLHiS0Lid/17zkB4wddmsDPsg=="},
		PasswordHashTask{Password: "123456", Salt: "oKuNU/4X5bdWsEu4OnZlfoPgEPIk8g=="},
		PasswordHashTask{Password: "9KCe4bFP7KxEqrDCzK", Salt: "0tarr2rbWgsjiPzxVIHSeKmuugUUhQ=="},
	}

	totalJobs := int64(len(tasks))
	p := pool.NewPool(initialWorkers, len(tasks))
	p.Start(ctx)

	go func() {
		for progress := range p.Progress() {
			bar := drawProgressBar(progress.Completed+progress.Failed, totalJobs)

			fmt.Printf("\rTelemetry: %s | Completed: %d | Failed: %d | Total: %d ",
				bar, progress.Completed, progress.Failed, totalJobs)
		}

		fmt.Println()
	}()

	go func() {
		for i, t := range tasks {
			p.Submit(pool.Job{
				ID:   i + 1,
				Task: t,
			})

			time.Sleep(200 * time.Millisecond)

			if i == 2 {
				fmt.Println("[Main] Scaling UP: Spawning +2 workers...")
				p.ScaleUp(ctx, 2)
			}

			if i == 5 {
				fmt.Println("\n[Main] Scaling DOWN: Terminating -2 workers...")
				p.ScaleDown(2)
			}
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
