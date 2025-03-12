package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var (
	words     = flag.Int("words", 2, "No. words in the petname")
	separator = flag.String("separator", " ", "Separator between words")
)

func main() {
	// Initialize Redis client
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	channel := "NewChannel"

	// Parse command-line flags
	flag.Parse()

	log.Printf("Starting publisher... Publishing to channel: %s\n", channel)

	for i := 1; i <= 5; i++ {
		name := petname.Generate(*words, *separator)
		message := fmt.Sprintf("Message %s", name)

		log.Printf("Publishing message %d: %s\n", i, message)

		err := client.Publish(ctx, channel, message).Err()
		if err != nil {
			log.Fatalf("Could not publish message: %v", err)
		} else {
			log.Printf("Successfully published message %d\n", i)
		}

		time.Sleep(4 * time.Second) // Simulate delay between messages
	}

	log.Println("Finished publishing messages.")
}
