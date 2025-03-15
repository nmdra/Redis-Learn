package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
	geoKey      = "drivers"
)

func main() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	fmt.Println("🚖 Ride-Sharing App using Redis Geo")

	// Menu
	for {
		fmt.Println("\n1. Add Driver")
		fmt.Println("2. Find Nearby Drivers")
		fmt.Println("3. Calculate Distance")
		fmt.Println("4. Exit")

		fmt.Print("Choose an option: ")
		choice := readInput()

		switch choice {
		case "1":
			addDriver()
		case "2":
			findNearbyDrivers()
		case "3":
			calculateDistance()
		case "4":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice! Try again.")
		}
	}
}

// Read user input
func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return input[:len(input)-1]
}

// Add a driver's location
func addDriver() {
	fmt.Print("Enter Driver Name: ")
	driver := readInput()

	fmt.Print("Enter Latitude: ")
	lat, _ := strconv.ParseFloat(readInput(), 64)

	fmt.Print("Enter Longitude: ")
	lon, _ := strconv.ParseFloat(readInput(), 64)

	_, err := redisClient.GeoAdd(ctx, geoKey, &redis.GeoLocation{
		Name:      driver,
		Longitude: lon,
		Latitude:  lat,
	}).Result()
	if err != nil {
		log.Fatalf("Failed to add driver: %v", err)
	}

	fmt.Printf("✅ Driver %s added at (%f, %f)\n", driver, lat, lon)
}

// Find nearby drivers within a given radius
func findNearbyDrivers() {
	fmt.Print("Enter Rider Latitude: ")
	lat, _ := strconv.ParseFloat(readInput(), 64)

	fmt.Print("Enter Rider Longitude: ")
	lon, _ := strconv.ParseFloat(readInput(), 64)

	fmt.Print("Enter Search Radius (km): ")
	radius, _ := strconv.ParseFloat(readInput(), 64)

	drivers, err := redisClient.GeoSearch(ctx, geoKey, &redis.GeoSearchQuery{
		Longitude:  lon,
		Latitude:   lat,
		Radius:     radius,
		RadiusUnit: "km",
		Sort:       "ASC", // Closest first
	}).Result()
	if err != nil {
		log.Fatalf("Failed to search drivers: %v", err)
	}

	if len(drivers) == 0 {
		fmt.Println("❌ No drivers found nearby.")
		return
	}

	fmt.Println("🚕 Nearby Drivers:")
	for i, driver := range drivers {
		fmt.Printf("%d. %s\n", i+1, driver)
	}
}

// Calculate distance between a rider and a driver
func calculateDistance() {
	fmt.Print("Enter First Location Name: ")
	loc1 := readInput()

	fmt.Print("Enter Second Location Name: ")
	loc2 := readInput()

	dist, err := redisClient.GeoDist(ctx, geoKey, loc1, loc2, "km").Result()
	if err != nil {
		log.Fatalf("Failed to calculate distance: %v", err)
	}

	fmt.Printf("📏 Distance between %s and %s: %.2f km\n", loc1, loc2, dist)
}
