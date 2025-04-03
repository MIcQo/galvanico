package server

import (
	"galvanico/internal/auth"
	"galvanico/internal/broker"
	"galvanico/internal/config"
	"galvanico/internal/database"
	"github.com/rs/zerolog"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog/log"
)

const (
	idleTimeout = 5 * time.Second
)

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

	printRegisteredRoutes(app)

	return app
}

func printRegisteredRoutes(app *fiber.App) {
	if log.Logger.GetLevel() >= zerolog.DebugLevel {
		return
	}

	var routes = app.GetRoutes()

	for _, route := range routes {
		if route.Method == fiber.MethodHead {
			continue
		}
		if route.Method != fiber.MethodGet && route.Path == "/" {
			continue
		}
		log.Debug().
			Str("method", route.Method).
			Str("path", route.Path).
			Str("name", route.Name).
			Msg("registered route")
	}
}

func registerUnauthorizedRoutes(app *fiber.App, cfg *config.Config) {
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessEndpoint:  "/livez",
		ReadinessEndpoint: "/readyz",
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

	app.Post("/auth/login", auth.LoginHandler)
}

func registerAuthorizedRoutes(app *fiber.App, cfg *config.Config) {
	var key = cfg.Auth.GetJWTKey()

	app.Use(jwtware.New(jwtware.Config{
		ErrorHandler: auth.ErrorHandler,
		SigningKey:   jwtware.SigningKey{Key: key},
	}))

	app.Get("/authorized", func(c *fiber.Ctx) error {
		var usr, err = auth.GetUser(c)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"user": usr,
		})
	})
}
