package server

import (
	"galvanico/internal/broker"
	"galvanico/internal/config"
	"galvanico/internal/database"
	"github.com/ansrivas/fiberprometheus/v2"
	"time"

	"github.com/goccy/go-json"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog/log"
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
		JSONDecoder: json.Unmarshal,
		JSONEncoder: json.Marshal,
	})

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log.Logger,
	}))

	app.Use(requestid.New())
	app.Use(recover.New())

	registerUnauthorizedRoutes(app, cfg)
	registerAuthorizedRoutes(app, cfg)

	return app
}

func registerUnauthorizedRoutes(app *fiber.App, cfg *config.Config) {
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(_ *fiber.Ctx) bool {
			return true
		},
		ReadinessProbe: func(ctx *fiber.Ctx) bool {
			return database.Connection().PingContext(ctx.Context()) == nil && broker.Connection().IsConnected()
		},
	}))

	prometheus := fiberprometheus.NewWithDefaultRegistry(cfg.AppName)
	prometheus.RegisterAt(app, "/metrics")
	prometheus.SetSkipPaths([]string{"/ping", "/readyz", "/livez"})
	app.Use(prometheus.Middleware)

}

func registerAuthorizedRoutes(_ *fiber.App, _ *config.Config) {

}
