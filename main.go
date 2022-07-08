package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var ctx = context.Background()

func main() {
	app := fiber.New()
	client := connectDatabase()

	app.Use(recover.New())

	app.Get("/:flavour/:version", func(c *fiber.Ctx) error {

		dbstate := checkDatabase(client)

		if(!dbstate) {
			return fiber.ErrServiceUnavailable
		}

		flavour := c.Params("flavour")
		version := c.Params("version")

		exists, err := client.HExists(ctx, flavour, version).Result()

		if !exists {
			return fiber.ErrNotFound
		}

		if err != nil {
			return fiber.ErrInternalServerError
		}

		resMap, err := client.HGet(ctx, flavour, version).Result()

		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.Redirect(resMap)
	})

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}

func connectDatabase() *redis.Client{
	client := redis.NewClient(&redis.Options{
        Addr:     net.JoinHostPort(os.Getenv("REDISHOST"),os.Getenv("REDISPORT")),
		Password: os.Getenv("REDISPASSWORD"),
		Username: os.Getenv("REDISUSER"),
    })
	return client
}

func checkDatabase(client *redis.Client) bool {
	_, err := client.Ping(ctx).Result()

	if(err != nil){
		return false
	} else {
		return true
	}
}