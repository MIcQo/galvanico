package server

import (
	"errors"
	"galvanico/internal/auth"
	"galvanico/internal/broker"
	"galvanico/internal/config"
	"galvanico/internal/database"
	"galvanico/internal/game/city"
	"galvanico/internal/game/user"
	"sync"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzerolog"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	idleTimeout = 5 * time.Second
)

var (
	lastReadinessCheck     time.Time
	lastReadinessResult    bool
	readinessCheckMutex    sync.Mutex
	readinessCheckInterval = 10 * time.Second
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

	app.Use(cors.New(cors.Config{
		// in the future, replace with http://localhost:5173, http://localhost:8080, https://galvanico.io
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false, // need to switch to true for production
	}))

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log.Logger,
	}))

	app.Use(requestid.New())
	app.Use(recover.New())

	app.Use(healthcheck.New(healthcheck.Config{
		LivenessEndpoint:  "/livez",
		ReadinessEndpoint: "/readyz",
		LivenessProbe: func(_ *fiber.Ctx) bool {
			return true
		},
		ReadinessProbe: func(ctx *fiber.Ctx) bool {
			readinessCheckMutex.Lock()
			defer readinessCheckMutex.Unlock()

			now := time.Now()
			if now.Sub(lastReadinessCheck) < readinessCheckInterval {
				return lastReadinessResult
			}

			result := database.Connection().PingContext(ctx.Context()) == nil && broker.Connection().IsConnected()
			lastReadinessCheck = now
			lastReadinessResult = result
			return result
		},
	}))

	prometheus := fiberprometheus.NewWithDefaultRegistry(cfg.AppName)
	prometheus.RegisterAt(app, "/metrics")
	prometheus.SetSkipPaths([]string{"/ping", "/readyz", "/livez"})
	app.Use(prometheus.Middleware)

	registerAPIroutes(app, cfg)
	registerStaticRoutes(app)

	printRegisteredRoutes(app)

	return app
}

func registerStaticRoutes(app *fiber.App) {
	app.Static("/", "./public")

	app.Get("/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./public/index.html")
	})
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

func registerAPIroutes(app *fiber.App, cfg *config.Config) {
	var g = app.Group("")
	{
		registerUnauthorizedRoutes(g, cfg)
		registerAuthorizedRoutes(g, cfg)
	}
}

func registerUnauthorizedRoutes(app fiber.Router, cfg *config.Config) {
	var userRepo = user.NewUserRepository(database.Connection())
	var publisher = broker.NewNatsPublisher(broker.Connection())
	var userHandler = user.NewHandler(userRepo, user.NewService(userRepo, publisher), cfg)

	var ag = app.Group("/auth")
	{
		ag.Post("/login", userHandler.LoginHandler)
		ag.Post("/register", userHandler.RegisterHandler)
	}
}

func registerAuthorizedRoutes(app fiber.Router, cfg *config.Config) {
	var key = cfg.Auth.GetJWTKey()
	var userRepo = user.NewUserRepository(database.Connection())
	var cityRepo = city.NewRepository(database.Connection())
	var publisher = broker.NewNatsPublisher(broker.Connection())
	var userService = user.NewService(userRepo, publisher)
	var userHandler = user.NewHandler(userRepo, userService, cfg)
	var cityHandler = city.NewHandler(cityRepo, city.NewService(cityRepo), userService)
	var userMiddleware = user.NewUserMiddleware(userRepo)

	var api = app.Group("/api")
	{
		api.Use(jwtware.New(jwtware.Config{
			ErrorHandler: auth.ErrorHandler,
			SigningKey:   jwtware.SigningKey{Key: key},
		}))

		api.Use(userMiddleware.CheckNotBanned())

		var usr = api.Group("/user")
		{
			usr.Get("", userHandler.GetHandler)
			usr.Patch("/username", userHandler.ChangeUsernameHandler)
			usr.Patch("/password", userHandler.ChangePasswordHandler)
			var ct = usr.Group("/city")
			{
				ct.Get("", cityHandler.HandleGetUserCities)

				var cg = ct.Group("/:city")
				{
					cg.Get("/building/:slot/available", cityHandler.HandleAvailableSlotBuildings)
				}
			}
		}
	}
}
