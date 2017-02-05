package main

import (
	"benchmark"
	"benchmark/connection"
	"benchmark/report"
	"benchmark/suite"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of Tile38-benchmark:")
		flag.PrintDefaults()
	}

	connectionString := flag.String("connection", "127.0.0.1:9851", "Connection string to Tile38 server")
	requests := flag.Uint64("requests", 100000, "Number of requires will be send totall for each benchmark")
	clients := flag.Int("clients", 50, "Number of clients will be connected for each benchmark")
	keepalive := flag.Bool("keepalive", true, "Keep connection alive")
	flush := flag.Bool("flush", false, "Flush DB before running benchmark")
	flag.Parse()

	if *requests < 1 {
		fmt.Println("You should set at least 1 request")
		os.Exit(1)
	}

	if *clients < 1 {
		fmt.Println("You should set at least 1 client")
		os.Exit(1)
	}

	runners := make(map[string]suite.Runner)

	// Get command
	runners["GET(OBJECT)"] = &suite.GetCommand{suite.Object}
	runners["GET(POINT)"] = &suite.GetCommand{suite.Point}
	runners["GET(BOUNDS)"] = &suite.GetCommand{suite.Bounds}
	runners["GET(HASH)"] = &suite.GetCommand{suite.Geohash}

	// Set command
	runners["SET(POINT)"] = &suite.SetCommand{suite.Point}

	conn := connection.NewTCPConnection(*connectionString)
	reporter := report.NewCLIReporter()

	options := benchmark.BenchmarkOptions{
		Connection: conn,
		Reporter: reporter,
		Requests: *requests,
		Clients: *clients,
		Keepalive: *keepalive,
		Flush: *flush,
	}

	benchmark.Run(runners, options)
}