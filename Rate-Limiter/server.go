package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Global variables
var (
	ctx         = context.Background()
	redisClient *redis.Client
)

const (
	requestLimit = 5           // Max requests per minute
	timeWindow   = time.Minute // Rate limit time window
)

func main() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	// Initialize Gin router
	router := gin.Default()
	router.Static("/static", "./static") // Serve frontend files
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// API endpoint with rate limiting
	router.GET("/api/request", rateLimitMiddleware, handleRequest)

	// Start server
	fmt.Println("Server running on http://localhost:8089")
	router.Run(":8089")
}

// Middleware for rate limiting
func rateLimitMiddleware(c *gin.Context) {
	userID := c.ClientIP() // Use IP as the user identifier

	allowed, remaining := isAllowed(userID)
	if !allowed {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"message": "⛔ Rate limit exceeded. Try again later.",
		})
		c.Abort()
		return
	}

	c.Set("remaining", remaining) // Pass remaining requests to handler
	c.Next()
}

// Check if user request is allowed
func isAllowed(userID string) (bool, int64) {
	key := fmt.Sprintf("rate_limit:%s", userID) // Redis key per user

	// Get current request count
	count, err := redisClient.Get(ctx, key).Int64()
	if err != nil && err != redis.Nil {
		log.Fatalf("Redis error: %v", err)
	}

	if count >= requestLimit {
		return false, 0 // Block request
	}

	// Increment the counter
	newCount, err := redisClient.Incr(ctx, key).Result()
	if err != nil {
		log.Fatalf("Failed to increment counter: %v", err)
	}

	// Set expiration if first request
	if newCount == 1 {
		redisClient.Expire(ctx, key, timeWindow)
	}

	return true, requestLimit - newCount // Allowed request
}

// Handle API request
func handleRequest(c *gin.Context) {
	remaining := c.GetInt64("remaining")
	c.JSON(http.StatusOK, gin.H{
		"message":   "✅ Request successful!",
		"remaining": remaining,
	})
}
