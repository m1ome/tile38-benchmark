package report

import (
	"strings"
	"fmt"
	"time"
)

type CLIReporter struct {

}

func NewCLIReporter() *CLIReporter {
	return &CLIReporter{}
}

func (r *CLIReporter) AddReport(options ReporterOptions) error {
	// Heading printing
	fmt.Printf("====== %s ======\n", strings.ToUpper(options.Name))
	fmt.Printf("  %d requests completed in %s\n", options.Requests, options.Elapsed)
	fmt.Printf("  %d parallel clients\n", options.Clients)
	fmt.Printf("  keep alive: %t\n", options.Keepalive)
	fmt.Println()

	// Printing timings for requests range
	medians := getMedians(options.Timings, options.Requests)
	if len(medians) > 0 {
		for _, r := range medians {
			if r.p < 1.0 {
				fmt.Printf("<1%%(%d) `<=` %d milliseconds\n", r.r, r.t)
			} else {
				fmt.Printf("%.2f%% `<=` %d milliseconds\n", r.p, r.t)
			}
		}

		fmt.Println()
	}

	// Printing footer
	rps := float64(options.Requests)  / (float64(options.Elapsed) / float64(time.Second))
	fmt.Printf("%.2f requests per second\n", rps)
	fmt.Println()

	return nil
}

func (r *CLIReporter) Footer() error {
	fmt.Println("====== END ======")
	fmt.Println("Benchmark finished")
	fmt.Println()

	return nil
}
