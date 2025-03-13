package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	defer client.Close()

	queue := "alarm_queue"

	for {
		// Blocking pop (BLPOP) to wait for messages
		result, err := client.BLPop(ctx, 0, queue).Result()
		if err != nil {
			log.Fatalf("Error retrieving message: %v", err)
		}

		messageData := result[1]
		fmt.Println("Received message from queue:", messageData)

		// Parse JSON metadata
		var messageMetadata map[string]interface{}
		err = json.Unmarshal([]byte(messageData), &messageMetadata)
		if err != nil {
			log.Printf("Failed to parse metadata: %v", err)
			continue
		}

		messageKey := fmt.Sprintf("message:%s", messageMetadata["id"])
		messageDetails, err := client.HGetAll(ctx, messageKey).Result()
		if err != nil {
			log.Printf("Failed to fetch message details: %v", err)
			continue
		}

		if len(messageDetails) == 0 {
			log.Println("Message details not found (possibly expired).")
			continue
		}

		fmt.Println("Processing message:", messageDetails)

		time.Sleep(2 * time.Second)
	}
}
