package main

import (
	"context"
	"log"
	"os"

	"example.com/m/v2/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var ctx = context.Background()

func main() {
	app := fiber.New()
	client := database.Connect()

	app.Use(recover.New())
	app.Static("/", "./public")

	app.Get("/:flavour/:version", func(c *fiber.Ctx) error {

		dbstate := database.Check(client)

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