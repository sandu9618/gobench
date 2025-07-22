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
)

var rootCmd = &cobra.Command{
	Use:   "gobench",
	Short: "GoBench is a simple HTTP benchmark tool",
	Long: `GoBench is a tiny HTTP benchmarking tool written in Go.

Usage:
  gobench [flags]

Examples:
  gobench -u https://example.com -n 100 -c 5

Flags:
  -u, --url           URL to benchmark
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

		bench.RunBenchMark(url, totalReqs, concurrency)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL to benchmark")
	rootCmd.Flags().IntVarP(&totalReqs, "requests", "n", 1, "Total number of requests")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrent workers")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
