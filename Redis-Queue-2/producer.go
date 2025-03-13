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

	messageMetadata := map[string]interface{}{
		"id":             "msg123",
		"sender_id":      "deviceA",
		"signal_code":    "12345", // From a predefined list
		"criteria_index": 1,
	}

	metadataJSON, err := json.Marshal(messageMetadata)
	if err != nil {
		log.Fatalf("Failed to marshal metadata: %v", err)
	}

	queue := "alarm_queue"
	err = client.LPush(ctx, queue, metadataJSON).Err()
	if err != nil {
		log.Fatalf("Failed to push metadata to queue: %v", err)
	}

	ttl := 1 * time.Second
	messageKey := fmt.Sprintf("message:%s", messageMetadata["id"])

	err = client.HSet(ctx, messageKey, map[string]interface{}{
		"state":       "new",
		"phase":       "initial",
		"destination": "center1",
		"ttl":         ttl.Seconds(),
	}).Err()
	if err != nil {
		log.Fatalf("Failed to store message details: %v", err)
	}

	err = client.Expire(ctx, messageKey, ttl).Err()
	if err != nil {
		log.Fatalf("Failed to set expiration: %v", err)
	}

	fmt.Println("Message metadata pushed to the queue, full message details stored in Redis, and expiration set.")
}
