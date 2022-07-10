package database

import (
	"context"
	"net"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// Connect returns a redis client
func Connect() *redis.Client{
	client := redis.NewClient(&redis.Options{
        Addr:     net.JoinHostPort(os.Getenv("REDISHOST"),os.Getenv("REDISPORT")),
		Password: os.Getenv("REDISPASSWORD"),
		Username: os.Getenv("REDISUSER"),
    })
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