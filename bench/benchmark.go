package bench

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

type Result struct {
	Duration     time.Duration
	Success      bool
	StatusCode   int
	Error        string
	ResponseBody string
}

func worker(id int, url, method, body, contentType string, jobs <-chan int, results chan<- Result, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	defer wg.Done()
	for range jobs {
		start := time.Now()

		// Create request body if provided
		var reqBody io.Reader
		if body != "" {
			reqBody = strings.NewReader(body)
		}

		// Create request with specified method and body
		req, err := http.NewRequest(method, url, reqBody)
		if err != nil {
			results <- Result{
				Duration: time.Since(start),
				Success:  false,
				Error:    fmt.Sprintf("Failed to create request: %v", err),
			}
			continue
		}

		// Set content type header if provided
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}

		// Send request
		client := &http.Client{}
		resp, err := client.Do(req)
		duration := time.Since(start)

		if err != nil {
			results <- Result{
				Duration: duration,
				Success:  false,
				Error:    fmt.Sprintf("Request failed: %v", err),
			}
			continue
		}

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			results <- Result{
				Duration:   duration,
				Success:    false,
				StatusCode: resp.StatusCode,
				Error:      fmt.Sprintf("Failed to read response body: %v", err),
			}
			continue
		}

		// Check if status code indicates failure
		if resp.StatusCode >= 400 {
			results <- Result{
				Duration:     duration,
				Success:      false,
				StatusCode:   resp.StatusCode,
				Error:        fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status),
				ResponseBody: string(respBody),
			}
			continue
		}

		results <- Result{
			Duration:     duration,
			Success:      true,
			StatusCode:   resp.StatusCode,
			ResponseBody: string(respBody),
		}

		bar.Add(1)
	}
}

func RunBenchMark(url, method, body, contentType string, showResponse bool, totalReqs, concurrency int) {
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
		go worker(w, url, method, body, contentType, jobs, results, &wg, bar)
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
	var statusCodes = make(map[int]int)
	var errors = make(map[string]int)
	var sampleResponse string

	for res := range results {
		if res.Success {
			successCount++
			// Store first successful response as sample
			if sampleResponse == "" && res.ResponseBody != "" {
				sampleResponse = res.ResponseBody
			}
		} else {
			failCount++
			// Count status codes for failed requests
			if res.StatusCode > 0 {
				statusCodes[res.StatusCode]++
			}
			// Count error types
			if res.Error != "" {
				errors[res.Error]++
			}
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
	fmt.Println("URL :", url)
	fmt.Println("Method :", method)
	if body != "" {
		fmt.Println("Body :", body)
	}
	if contentType != "" {
		fmt.Println("Content-Type :", contentType)
	}
	fmt.Println("Total Requests :", totalReqs)
	fmt.Println("Success :", successCount)
	fmt.Println("Failed :", failCount)

	// Show failure details if there are any
	if failCount > 0 {
		fmt.Println("")
		fmt.Println("--- Failure Details ---")

		// Show status code breakdown
		if len(statusCodes) > 0 {
			fmt.Println("Status Codes:")
			for status, count := range statusCodes {
				fmt.Printf("  %d: %d requests\n", status, count)
			}
		}

		// Show error breakdown
		if len(errors) > 0 {
			fmt.Println("Errors:")
			for err, count := range errors {
				fmt.Printf("  %s: %d requests\n", err, count)
			}
		}
		fmt.Println("")
	}

	// Show sample response if requested and available
	if showResponse && sampleResponse != "" {
		fmt.Println("--- Sample Response ---")
		// Truncate long responses for readability
		if len(sampleResponse) > 500 {
			fmt.Println(sampleResponse[:500] + "...")
			fmt.Println("(Response truncated - first 500 characters shown)")
		} else {
			fmt.Println(sampleResponse)
		}
		fmt.Println("")
	}

	fmt.Println("Min Time :", minTime)
	fmt.Println("Max Time :", maxTime)
	fmt.Println("Avg Time :", avgTime)
	fmt.Println("Total Duration:", totalTime)
	fmt.Printf("Requests/sec : %.2f\n", float64(totalReqs)/totalTime.Seconds())
}
