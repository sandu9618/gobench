package cmd

import (
	"fmt"
	"os"

	"github.com/sandu9618/gobench/bench"
	"github.com/spf13/cobra"
)

var (
	url         string
	totalReqs   int
	concurrency int
	method      string
	body        string
	contentType string
)

var rootCmd = &cobra.Command{
	Use:   "gobench",
	Short: "GoBench is a simple HTTP benchmark tool",
	Long: `GoBench is a tiny HTTP benchmarking tool written in Go.

Usage:
  gobench [flags]

Examples:
  gobench -u https://example.com -n 100 -c 5
  gobench -u https://api.example.com/users -m POST -n 50 -c 10
  gobench -u https://api.example.com/users -m POST -d '{"name":"John","email":"john@example.com"}' -H "application/json" -n 50 -c 10

Flags:
  -u, --url           URL to benchmark
  -m, --method        HTTP method to use (GET, POST, PUT, DELETE, etc.)
  -d, --data          Request body data (for POST/PUT requests)
  -H, --header        Content-Type header (default: application/json for POST/PUT with data)
  -n, --requests      Number of requests to send
  -c, --concurrency   Number of concurrent workers
`,
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			cmd.Help()
			return
		}

		if totalReqs <= 0 || concurrency <= 0 {
			fmt.Println("Number of requests and concurrency must be greater than zero")
			os.Exit(1)
		}

		// Validate HTTP method
		if method == "" {
			method = "GET" // Default to GET if not specified
		}

		// Set default content type for POST/PUT with data
		if body != "" && contentType == "" {
			if method == "POST" || method == "PUT" {
				contentType = "application/json"
			}
		}

		bench.RunBenchMark(url, method, body, contentType, totalReqs, concurrency)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL to benchmark")
	rootCmd.Flags().StringVarP(&method, "method", "m", "GET", "HTTP method to use (GET, POST, PUT, DELETE, etc.)")
	rootCmd.Flags().StringVarP(&body, "data", "d", "", "Request body data (for POST/PUT requests)")
	rootCmd.Flags().StringVarP(&contentType, "header", "H", "", "Content-Type header")
	rootCmd.Flags().IntVarP(&totalReqs, "requests", "n", 1, "Total number of requests")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrent workers")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
