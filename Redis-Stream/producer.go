package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ctx        = context.Background()
	streamName = "mystream"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	for i := 1; i <= 15; i++ {
		msgID, err := client.XAdd(ctx, &redis.XAddArgs{
			Stream: streamName,
			Values: map[string]interface{}{
				"event": fmt.Sprintf("message-%d", i),
			},
		}).Result()

		if err != nil {
			log.Printf("Failed to add message: %v", err)
		} else {
			fmt.Printf("Produced message ID: %s\n", msgID)
		}

		time.Sleep(1 * time.Second) // Simulating interval between messages
	}
}
