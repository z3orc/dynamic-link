package main

import (
	"context"
	"log"
	"os"
	"strings"

	"example.com/m/v2/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Initialize various "global" variables
var (
	ctx = context.Background()
	projectName = "Dynamic Link"
)

func main() {
	// Initialize periodic database sync
	go database.PeriodicSync()

	// Initialize server
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: false,
		StrictRouting: false,
		GETOnly: true,
		DisableKeepalive: true,
		ServerHeader:  projectName,
		AppName: projectName,
	})
	
	// Initialize database client
	client := database.Connect()

	// Initialize middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(compress.New())
	app.Static("/", "./public")

	// Main route, used to fetch download url and redirect to the correct endpoint
	app.Get("/:flavour/:version", func(c *fiber.Ctx) error {

		dbstate := database.Check(client)

		if(!dbstate) {
			return fiber.ErrServiceUnavailable
		}

		flavour := strings.ToLower(c.Params("flavour"))
		version := strings.ToLower(c.Params("version"))

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

	// Starting server + logging fatal errors
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))

}