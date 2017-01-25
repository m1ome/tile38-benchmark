package main

import (
	"gopkg.in/redis.v5"
	"testing"
)

func connection() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:9851",
	})
}

func benchmarkNearby(b *testing.B, postfix string) {
	connection := connection()
	b.ResetTimer()

	for n :=0; n < b.N; n++ {
		cmd := redis.NewStringCmd("NEARBY", "benchmark-" + postfix, "POINT", "100.001", "100.001", "1000")
		connection.Process(cmd)
	}

	connection.Close()
}

func benchmarkNearbyWithLimitAndDistance(b *testing.B, postfix string) {
	connection := connection()
	b.ResetTimer()

	for n :=0; n < b.N; n++ {
		cmd := redis.NewStringCmd("NEARBY", "benchmark-" + postfix, "LIMIT", 100, "DISTANCE", "POINT", "100.001", "100.001", "1000")
		connection.Process(cmd)
	}

	connection.Close()
}

func benchmarkNearbyWithLimit(b *testing.B, postfix string) {
	connection := connection()
	b.ResetTimer()

	for n :=0; n < b.N; n++ {
		cmd := redis.NewStringCmd("NEARBY", "benchmark-" + postfix, "LIMIT", 100, "POINT", "100.001", "100.001", "1000")
		connection.Process(cmd)
	}

	connection.Close()
}

// 1 Million records
func BenchmarkNearby1M(b *testing.B) {
	benchmarkNearby(b, "1m")
}

func BenchmarkNearbyWithLimit1M(b *testing.B) {
	benchmarkNearbyWithLimit(b, "1m")
}

func BenchmarkNearbyWithLimitAndDistance1M(b *testing.B) {
	benchmarkNearbyWithLimitAndDistance(b, "1m")
}

// 3 Million records
func BenchmarkNearby3M(b *testing.B) {
	benchmarkNearby(b, "3m")
}

func BenchmarkNearbyWithLimit3M(b *testing.B) {
	benchmarkNearbyWithLimit(b, "3m")
}

func BenchmarkNearbyWithLimitAndDistance3M(b *testing.B) {
	benchmarkNearbyWithLimitAndDistance(b, "3m")
}

// 5 Million records
func BenchmarkNearby5M(b *testing.B) {
	benchmarkNearby(b, "5m")
}

func BenchmarkNearbyWithLimit5M(b *testing.B) {
	benchmarkNearbyWithLimit(b, "5m")
}

func BenchmarkNearbyWithLimitAndDistance5M(b *testing.B) {
	benchmarkNearbyWithLimitAndDistance(b, "5m")
}