package database

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// Connects & returns a redis client
func Connect() (*redis.Client){
	url, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Println("Could not parse url")
	}

	client := redis.NewClient(url)

	return client
}

// Check the state of a redis client
func Check(client *redis.Client) bool {
	_, err := client.Ping(ctx).Result()

	if(err != nil){
		return false
	} else {
		return true
	}
}