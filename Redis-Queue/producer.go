package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
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

	// Push tasks to queue
	tasks := []string{"task1", "task2", "task3", "task4", "task5", "task6"}

	for _, task := range tasks {
		err := client.RPush(ctx, queue, task).Err()
		if err != nil {
			log.Fatalf("Failed to push task: %v", err)
		}
		fmt.Println("Added:", task)

		// Sleep for a random duration (0-5 seconds)
		delay := time.Duration(rand.Intn(5)) * time.Second
		time.Sleep(delay)
	}
}
