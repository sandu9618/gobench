package bench

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

type Result struct {
	Duration time.Duration
	Success  bool
}

func worker(id int, url string, jobs <-chan int, results chan<- Result, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	defer wg.Done()
	for range jobs {
		start := time.Now()
		resp, err := http.Get(url)
		duration := time.Since(start)

		if err != nil || resp.StatusCode >= 400 {
			results <- Result{Duration: duration, Success: false}
			continue
		}

		results <- Result{Duration: duration, Success: true}
		resp.Body.Close()

		bar.Add(1)
	}
}

func RunBenchMark(url string, totalReqs, concurrency int) {
	jobs := make(chan int, totalReqs)
	results := make(chan Result, totalReqs)
	var wg sync.WaitGroup

	bar := progressbar.NewOptions(totalReqs,
		progressbar.OptionSetDescription("Sending requests"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "#",
			SaucerHead:    ">",
			SaucerPadding: "-",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	for w := 1; w <= concurrency; w++ {
		wg.Add(1)
		go worker(w, url, jobs, results, &wg, bar)
	}

	for j := 0; j < totalReqs; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	close(results)

	var successCount, failCount int
	var totalTime time.Duration
	var minTime, maxTime time.Duration

	for res := range results {
		if res.Success {
			successCount++
		} else {
			failCount++
		}

		totalTime += res.Duration

		if minTime == 0 || res.Duration < minTime {
			minTime = res.Duration
		}
		if res.Duration > maxTime {
			maxTime = res.Duration
		}
	}

	avgTime := totalTime / time.Duration(totalReqs)

	fmt.Println("")
	fmt.Println("====== GoBench Result ======")
	fmt.Println("Total Requests :", totalReqs)
	fmt.Println("Success :", successCount)
	fmt.Println("Failed :", failCount)
	fmt.Println("Min Time :", minTime)
	fmt.Println("Max Time :", maxTime)
	fmt.Println("Avg Time :", avgTime)
	fmt.Println("Total Duration:", totalTime)
	fmt.Printf("Requests/sec : %.2f\n", float64(totalReqs)/totalTime.Seconds())
}
