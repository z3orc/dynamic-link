package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/z3orc/dynamic-link/database"
	"github.com/z3orc/dynamic-link/endpoints/paper"
	"github.com/z3orc/dynamic-link/endpoints/purpur"
	"github.com/z3orc/dynamic-link/endpoints/vanilla"
	"github.com/z3orc/dynamic-link/util"
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

	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUID,
	}))

	app.Static("/", "./public")

	// Main route, used to fetch download url and redirect to the correct endpoint
	app.Get("/:flavour/:version", func(c *fiber.Ctx) error {
		var result string
		var err error

		flavour := strings.ToLower(c.Params("flavour"))
		version := strings.ToLower(c.Params("version"))
		dbstate := database.Check(client)

		if flavour != "vanilla" && flavour != "paper" && flavour != "purpur"{
			return fiber.ErrBadRequest
		}

		fmt.Println(dbstate)

		switch dbstate {
		case true:
			exists, err := client.HExists(ctx, flavour, version).Result()

			if !exists {
				return fiber.ErrNotFound
			}

			if err != nil {
				return fiber.ErrInternalServerError
			}

			result, err = client.HGet(ctx, flavour, version).Result()

			if err != nil {
				return fiber.ErrInternalServerError
			}

		case false:
			switch flavour {
			case "vanilla":
				result, err = vanilla.GetDownloadUrl(version)
			case "paper":
				result, err = paper.GetDownloadUrl(version)
			case "purpur":
				result, err = purpur.GetDownloadUrl(version)
			}

			if err != nil && err.Error() == "404" {
				return fiber.ErrNotFound
			} else if err != nil{
				return fiber.ErrInternalServerError
			}
		}

		err = util.CheckUrl(result)
		if err != nil {
			return fiber.ErrServiceUnavailable
		}

		return c.Redirect(result)
	})

	// Starting server + logging fatal errors
	log.Fatal(app.Listen(":8080"))
}