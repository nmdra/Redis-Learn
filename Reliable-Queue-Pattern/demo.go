package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
)

// Adds a message to the pending queue
func enqueueMessage(message string) error {
	return rdb.LPush(ctx, "pending_queue", message).Err()
}

// Moves message from pending queue to processing queue
func dequeueMessage() (string, error) {
	msg, err := rdb.LMove(ctx, "pending_queue", "processing_queue", "RIGHT", "LEFT").Result()
	if err == redis.Nil {
		return "", nil // No messages available
	}
	return msg, err
}

// Process messages reliably
func processMessages() {
	for {
		message, err := dequeueMessage()
		if err != nil {
			log.Println("Error moving message:", err)
			continue
		}

		if message == "" {
			time.Sleep(1 * time.Second) // No messages, wait before polling
			continue
		}

		fmt.Println("Processing message:", message)

		time.Sleep(1 * time.Second) // Slow Down process
		// Simulating failure or success
		if time.Now().Unix()%2 == 0 {
			fmt.Println("Message processed successfully:", message)
			rdb.LRem(ctx, "processing_queue", 1, message) // Remove from processing queue
		} else {
			fmt.Println("Processing failed, retrying:", message)
			// Move back to pending queue
			rdb.LMove(ctx, "processing_queue", "pending_queue", "RIGHT", "LEFT")
		}
	}
}

func main() {
	// Simulating producer
	enqueueMessage("task-1")
	enqueueMessage("task-2")
	enqueueMessage("task-3")
	enqueueMessage("task-4")
	enqueueMessage("task-5")
	enqueueMessage("task-6")
	enqueueMessage("task-7")

	// Start processing
	processMessages()
}
