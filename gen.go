package main

import (
	"fmt"
	"math/rand"

	redis "gopkg.in/redis.v5"
)

const Million = 100000

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:9851",
	})

	defer client.Close()

	fmt.Println("Generate points for Tile38")

	generatePoints(client, Million, "1m")
	generatePoints(client, 3 * Million, "3m")
	generatePoints(client, 5 * Million, "5m")

	fmt.Println("Tile38 populated successfully")
}


func randPosition(minLat, minLon, maxLat, maxLon float64) (float64, float64) {
	lat, lon := (rand.Float64()*(maxLat-minLat))+minLat, (rand.Float64()*(maxLon-minLon))+minLon
	return lat, lon
}

func generatePoints(client *redis.Client, count int, postfix string) {
	fmt.Printf("We will generate %d points\n", count)

	for i := 0; i < count; i++ {
		lat, lon := randPosition(50.001, 50.001, 50.01, 50.01)
		cmd := redis.NewStringCmd("SET", "benchmark-" + postfix, i, "POINT", lat, lon)
		client.Process(cmd)

		fmt.Printf("Progress: %d/%d\n", i + 1, count)
	}
}