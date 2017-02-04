package main

import (
	"benchmark"
	"benchmark/connection"
	"benchmark/report"
	"benchmark/suite"
)

func main() {
	runners := make(map[string]suite.Runner)

	runners["GET"] = &suite.GetCommand{}
	runners["SET(POINT)"] = &suite.SetCommand{}

	conn := connection.NewTCPConnection("127.0.0.1:9851")
	reporter := report.NewCLIReporter()

	options := benchmark.BenchmarkOptions{
		Connection: conn,
		Reporter: reporter,
		Requests: 10000,
		Clients: 30,
		Keepalive: true,
	}

	benchmark.Run(runners, options)
}