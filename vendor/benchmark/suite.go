package benchmark

import (
	"time"
	"sync"
	"fmt"
	"sort"
)

type Suite struct {
	mu sync.Mutex
	done uint64
	options SuiteOptions
	startTime time.Time
	timing map[int]int
	errors uint64
	success uint64
}

type SuiteOptions struct {
	Connections int
	Requests uint64
	Keepalive bool
}

func NewSuite(opts SuiteOptions) *Suite {
	if opts.Connections == 0 {
		opts.Connections = 50
	}

	if opts.Requests == 0 {
		opts.Requests = 100000
	}

	return &Suite{options:opts, timing: make(map[int]int)}
}

func (s *Suite) Run() {
	var wg sync.WaitGroup
	wg.Add(s.options.Connections)

	s.startTime = time.Now()
	keepalive := s.options.Keepalive
	for i := 0; i < s.options.Connections; i++ {
		go func() {
			c := NewConnection(ConnectionOptions{})

			//if keepalive {
			//	c.Dial()
			//	defer c.Close()
			//}

			for {
				if !keepalive {
					c.Dial()
				}

				s.mu.Lock()
				if s.done >= s.options.Requests {
					s.mu.Unlock()
					break
				}
				s.mu.Unlock()

				t := time.Now()
				_, err := c.Do("GET fleet truck")
				e := time.Since(t)

				if err != nil {
					panic(err)
				}

				s.mu.Lock()
				r := int(e / time.Millisecond)
				if r == 0 || r == 1 {
					r = 1
				}
				if _, ok := s.timing[r]; !ok {
					s.timing[r] = 0
				}
				s.timing[r]++
				s.done++
				s.mu.Unlock()

				if !keepalive {
					c.Close()
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()

	elapsed := time.Since(s.startTime)

	// Main statistics
	fmt.Println("====== GET ======")
	fmt.Printf("  %d requests completed in %s\n", s.options.Requests, elapsed)
	fmt.Printf("  %d parallel clients\n", s.options.Connections)
	fmt.Println()

	// Percentage and medians
	var keys []int
	for i := range s.timing {
		keys = append(keys, i)
	}
	sort.Ints(keys)

	var summary float64
	for _, i := range keys {
		c := s.timing[i]
		p := float64(c) / float64(s.options.Requests) * 100
		summary += p

		if summary >= 100 {
			summary = 100.00
		}

		fmt.Printf("%.2f `<=` %d milliseconds\n", summary, i)

		if summary >= 100 {
			break
		}
	}

	rps := float64(s.options.Requests)  / (float64(elapsed) / float64(time.Second))
	fmt.Printf("%.2f requests per second", rps)
}