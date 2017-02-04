package report

import (
	"time"
	"sort"
)

type ReporterOptions struct {
	Name string
	Requests uint64
	Elapsed time.Duration
	Clients int
	Keepalive bool
	Timings map[int]int
}

type Reporter interface {
	AddReport(opts ReporterOptions) error
	Footer() error
}


type Median struct {
	p float64
	t int
	r int
}

func getMedians(timing map[int]int, requests uint64) []Median {
	// Sorted set of timings
	var sorted []Median

	// Percentage and medians
	var keys []int
	for i := range timing {
		keys = append(keys, i)
	}
	sort.Ints(keys)

	var summary float64
	for _, i := range keys {
		c := timing[i]
		p := float64(c) / float64(requests) * 100
		summary += p

		sorted = append(sorted, Median{p, i, c})

		if summary >= 100 {
			summary = 100.00
		}

		if summary >= 100 {
			break
		}
	}

	return sorted
}
