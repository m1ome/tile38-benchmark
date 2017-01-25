package main

import (
	"fmt"

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


func generatePoints(client *redis.Client, count int, postfix string) {
	fmt.Printf("We will generate %d points\n", count)

	for i := 0; i < count; i++ {
		cmd := redis.NewStringCmd("SET", "benchmark-" + postfix, i, "POINT", 100, 100)
		client.Process(cmd)

		fmt.Printf("Progress: %d/%d\n", i + 1, count)
	}
}