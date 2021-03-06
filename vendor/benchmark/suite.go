package benchmark

import (
	"sync"
	"benchmark/connection"
	"benchmark/report"
	"benchmark/suite"
	"time"
	"fmt"
	"os"
	"strings"
)

type BenchmarkOptions struct {
	Connection connection.Connector
	Reporter report.Reporter
	Keepalive bool
	Clients int
	Requests uint64
	Flush bool
}

func Run(tests map[string]suite.Runner, opts BenchmarkOptions) {
	conn, err := opts.Connection.Dial()
	defer conn.Close()

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	if opts.Flush {
		fmt.Println("Flushing database")

		err = conn.Write("FLUSHDB")
		if err != nil {
			fmt.Printf("Error(flush): %s", err.Error())
			os.Exit(1)
		}

		m, err := conn.Read()
		if err != nil {
			fmt.Printf("Error(flush): %s", err.Error())
			os.Exit(1)
		}

		if !strings.Contains(string(m), "\"ok\":true") {
			fmt.Printf("Error: %s", string(m))
			os.Exit(1)
		}
	}

	for name, runner := range tests {
		runner.Up(conn)

		options := command(runner, opts.Connection, opts.Keepalive, opts.Clients, opts.Requests)
		options.Name = name
		opts.Reporter.AddReport(options)

		runner.Down(conn)
	}

	opts.Reporter.Footer()
}

func command(runner suite.Runner, conn connection.Connector, keepalive bool, clients int, requests uint64) report.ReporterOptions {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var done uint64

	timing := make(map[int]int)
	wg.Add(clients)

	var elapsed time.Duration
	for i := 0; i < clients; i++ {
		go func() {
			var c connection.ConnectorReadWriter
			var err error

			if keepalive {
				c, err = conn.Dial()
				if err != nil {
					fmt.Printf("Connection error: %s\n", err)
					os.Exit(1)
				}
			}

			for {
				mu.Lock()
				if done >= requests {
					mu.Unlock()
					break
				}
				done++
				mu.Unlock()

				// Establish connection if not keepalive
				if !keepalive {
					c, err = conn.Dial()
					if err != nil {
						fmt.Printf("Connection error: %s\n", err)
						os.Exit(1)
					}
				}

				t := time.Now()
				err = runner.Fire(c)
				e := time.Since(t)

				mu.Lock()
				elapsed += e
				mu.Unlock()

				if err != nil {
					fmt.Printf("Fire error: %s\n", err.Error())
					c.Close()
					os.Exit(1)
				}

				err = runner.Match(c)
				if err != nil {
					fmt.Printf("Match error: %s\n", err.Error())
					c.Close()
					os.Exit(1)
				}

				mu.Lock()
				r := int(e / time.Millisecond)
				if r == 0  {
					r = 1
				}
				if _, ok := timing[r]; !ok {
					timing[r] = 0
				}
				timing[r]++
				mu.Unlock()

				if !keepalive {
					c.Close()
				}
			}

			if keepalive {
				c.Close()
			}

			wg.Done()
		}()
	}

	wg.Wait()
	return report.ReporterOptions{
		Elapsed: elapsed,
		Requests: requests,
		Clients: clients,
		Keepalive: keepalive,
		Timings: timing,
	}
}