package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer client.Close()

	queue := "task_queue"

	for {
		// Blocking pop (waits for new task)
		task, err := client.BLPop(ctx, 0, queue).Result()
		if err != nil {
			log.Fatalf("Error consuming task: %v", err)
		}
		fmt.Println("Consumed:", task[1])

		// Sleep for 1 second before consuming the next task
		time.Sleep(1 * time.Second)
	}
}
