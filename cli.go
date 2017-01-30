package main

import (
	"benchmark"
)

func main() {
	opts := benchmark.SuiteOptions{Keepalive: false}

	s := benchmark.NewSuite(opts)
	s.Run()

}