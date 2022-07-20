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

var (
	ctx = context.Background()
	projectName = "Dynamic Link"
	appConfig = fiber.Config{
		Prefork:       false,
		CaseSensitive: false,
		StrictRouting: false,
		GETOnly: true,
		DisableKeepalive: true,
		ServerHeader:  projectName,
		AppName: projectName,
	}
)

func main() {
	app := fiber.New(appConfig)
	client := database.Connect()

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(compress.New())
	app.Static("/", "./public")

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

	database.Sync()
	go database.HeartBeat()

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))

}