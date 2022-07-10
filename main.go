package main

import (
	"context"

	"example.com/m/v2/database"
)

var ctx = context.Background()
var projectName = "mc-dynamic-link"

func main() {

	database.SyncPaper()

	// app := fiber.New(fiber.Config{
	// 	Prefork:       false,
	// 	CaseSensitive: false,
	// 	StrictRouting: false,
	// 	GETOnly: true,
	// 	DisableKeepalive: true,
	// 	ServerHeader:  projectName,
	// 	AppName: projectName,
	// })
	// client := database.Connect()

	// app.Use(logger.New())
	// app.Use(recover.New())
	// app.Use(compress.New())
	// app.Static("/", "./public")

	// app.Get("/:flavour/:version", func(c *fiber.Ctx) error {

	// 	dbstate := database.Check(client)

	// 	if(!dbstate) {
	// 		return fiber.ErrServiceUnavailable
	// 	}

	// 	flavour := c.Params("flavour")
	// 	version := c.Params("version")

	// 	exists, err := client.HExists(ctx, flavour, version).Result()

	// 	if !exists {
	// 		return fiber.ErrNotFound
	// 	}

	// 	if err != nil {
	// 		return fiber.ErrInternalServerError
	// 	}

	// 	resMap, err := client.HGet(ctx, flavour, version).Result()

	// 	if err != nil {
	// 		return fiber.ErrInternalServerError
	// 	}

	// 	return c.Redirect(resMap)
	// })

	// log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}