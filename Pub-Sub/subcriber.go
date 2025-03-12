package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	channel := "NewChannel"
	log.Printf("Subscribing to channel: %s\n", channel)

	// Subscribe to the channel
	sub := client.Subscribe(ctx, channel)
	ch := sub.Channel()

	log.Println("Successfully subscribed! Waiting for messages...")

	// Listen for messages
	for msg := range ch {
		log.Printf("Received message from channel [%s]: %s\n", msg.Channel, msg.Payload)
	}
}
