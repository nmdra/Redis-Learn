package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	ctx            = context.Background()
	redisClient    *redis.Client
	leaderboardKey = "game_leaderboard"
)

func main() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	fmt.Println("🎮 Game Leaderboard using Redis Sorted Sets")

	// Menu
	for {
		fmt.Println("\n1. Add Player")
		fmt.Println("2. Update Score")
		fmt.Println("3. Show Leaderboard")
		fmt.Println("4. Get Player Rank")
		fmt.Println("5. Exit")

		fmt.Print("Choose an option: ")
		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			addPlayer()
		case "2":
			updateScore()
		case "3":
			showLeaderboard()
		case "4":
			getPlayerRank()
		case "5":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice! Try again.")
		}
	}
}

// Add a player with a score
func addPlayer() {
	fmt.Print("Enter Player Name: ")
	var player string
	fmt.Scanln(&player)

	fmt.Print("Enter Score: ")
	var score float64
	fmt.Scanln(&score)

	_, err := redisClient.ZAdd(ctx, leaderboardKey,
		redis.Z{Score: score, Member: player},
	).Result()
	if err != nil {
		log.Fatalf("Failed to add player: %v", err)
	}

	fmt.Printf("✅ Player %s added with score %.2f\n", player, score)
}

// Update player score
func updateScore() {
	fmt.Print("Enter Player Name: ")
	var player string
	fmt.Scanln(&player)

	fmt.Print("Enter Score to Add: ")
	var score float64
	fmt.Scanln(&score)

	_, err := redisClient.ZIncrBy(ctx, leaderboardKey, score, player).Result()
	if err != nil {
		log.Fatalf("Failed to update score: %v", err)
	}

	fmt.Printf("✅ Player %s score increased by %.2f\n", player, score)
}

// Show top players
func showLeaderboard() {
	fmt.Print("Enter number of top players to show: ")
	var limit int
	fmt.Scanln(&limit)

	players, err := redisClient.ZRevRangeWithScores(ctx, leaderboardKey, 0, int64(limit-1)).Result()
	if err != nil {
		log.Fatalf("Failed to fetch leaderboard: %v", err)
	}

	fmt.Println("🏆 Leaderboard:")
	for i, player := range players {
		fmt.Printf("%d. %s - %.2f\n", i+1, player.Member, player.Score)
	}
}

// Get player rank
func getPlayerRank() {
	fmt.Print("Enter Player Name: ")
	var player string
	fmt.Scanln(&player)

	rank, err := redisClient.ZRevRank(ctx, leaderboardKey, player).Result()
	if err != nil {
		log.Fatalf("Failed to get rank: %v", err)
	}

	fmt.Printf("🎖️ %s is ranked #%d\n", player, rank+1)
}
