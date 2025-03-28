package server

import (
	"galvanico/internal/config"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog/log"
	"time"
)

const idleTimeout = 5 * time.Second

func NewServer() *fiber.App {
	var cfg, err = config.Load()
	if err != nil {
		panic(err)
	}
	var app = fiber.New(fiber.Config{
		AppName:     cfg.AppName,
		IdleTimeout: idleTimeout,
	})

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log.Logger,
	}))

	app.Use(healthcheck.New())
	app.Use(requestid.New())
	app.Use(recover.New())

	return app
}
