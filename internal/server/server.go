package server

import (
	"errors"
	"galvanico/internal/auth"
	"galvanico/internal/broker"
	"galvanico/internal/config"
	"galvanico/internal/database"
	"galvanico/internal/game/user"
	"time"

	"github.com/rs/zerolog"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzerolog"
	jwtware "github.com/gofiber/contrib/jwt"
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

		// This is default fiber error handler, but we need to change to sending json
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Return status code with error message
			return c.Status(code).JSON(fiber.Map{"message": err.Error()})
		},
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

	var userHandler = user.NewHandler(user.NewUserRepository(database.Connection()))

	var ag = app.Group("/auth")
	{
		ag.Post("/login", userHandler.LoginHandler)
		ag.Post("/register", userHandler.RegisterHandler)
	}
}

func registerAuthorizedRoutes(app *fiber.App, cfg *config.Config) {
	var key = cfg.Auth.GetJWTKey()

	app.Use(jwtware.New(jwtware.Config{
		ErrorHandler: auth.ErrorHandler,
		SigningKey:   jwtware.SigningKey{Key: key},
	}))

	var userHandler = user.NewHandler(user.NewUserRepository(database.Connection()))

	var api = app.Group("/api")
	{
		var usr = api.Group("/user")
		{
			usr.Get("", userHandler.GetHandler)
		}
	}
}
