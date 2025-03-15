package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ctx        = context.Background()
	streamName = "mystream"
	groupName  = "mygroup"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	// Get consumer name from command-line arguments
	consumerName := "consumer1" // Default
	if len(os.Args) > 1 {
		consumerName = os.Args[1]
	}

	// Create Consumer Group (if not exists)
	err := client.XGroupCreateMkStream(ctx, streamName, groupName, "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Printf("Consumer group already exists or cannot be created: %v", err)
	}
	fmt.Printf("[%s] Started and waiting for messages...\n", consumerName)

	// Start consuming messages
	for {
		messages, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: consumerName,
			Streams:  []string{streamName, ">"},
			Count:    1,
			Block:    5 * time.Second,
		}).Result()

		if err == redis.Nil {
			continue // No new messages
		} else if err != nil {
			log.Printf("[%s] Failed to read messages: %v", consumerName, err)
			continue
		}

		// Process messages
		for _, stream := range messages {
			for _, msg := range stream.Messages {
				fmt.Printf("[%s] Received message: %v\n", consumerName, msg.Values)

				// Simulating message processing delay
				time.Sleep(2 * time.Second)

				// Acknowledge message
				_, err := client.XAck(ctx, streamName, groupName, msg.ID).Result()
				if err != nil {
					log.Printf("[%s] Failed to ACK message: %v", consumerName, err)
				} else {
					fmt.Printf("[%s] Acknowledged message ID: %s\n", consumerName, msg.ID)
				}
			}
		}
	}
}
