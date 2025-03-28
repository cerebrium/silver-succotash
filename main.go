package main

import (
	"log"
	"os"

	"hopdf.com/api/pdfcsv"
	"hopdf.com/api/weights_routes"
	"hopdf.com/db"
	"hopdf.com/localware"

	"hopdf.com/handlers/dashboard"
	"hopdf.com/handlers/index"
	"hopdf.com/handlers/notFound"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db_ref := db.LocalConnect()
	// Defer is a stack, this should
	// close the db connection at the
	// end.
	defer db_ref.Close()

	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	clerk.SetKey(clerkSecretKey)

	app := echo.New()

	// Middleware for whole app
	// Golang equivilant of helmet for node
	app.Use(middleware.Secure())

	// We don't want long running anything. If
	// we end up openeing sockets at some point
	// then we can reconsider.
	app.Use(middleware.Timeout())

	// Logger
	app.Use(middleware.Logger())

	// Allow panics not to crash the server
	app.Use(middleware.Recover())

	// Just start with a blocker for many requests
	app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))

	// Add db to all routes context
	app.Use(localware.AddDb(db_ref))

	// CORS
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:42069", "https://prate.pro/"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// adds middleware to gather metrics
	app.Use(echoprometheus.NewMiddleware("prate"))

	// Serve the htmx and other assets
	app.Static("/assets", "assets")

	// All attempts at routes that aren't
	// mounted at /v1 should take you to
	// the auth page. If Authenticated, it
	// will automatically re-route you to
	// the dashboard.
	app.GET("/*", func(c echo.Context) error {
		return index.IndexHandler(c)
	})

	// All routes mounted at /v1/* are authenticated users
	authApp := app.Group("/v1", localware.WithHeaderAuthorizationMiddleware)

	// Add user struct to context
	authApp.Use(localware.AddLocalUser)

	authApp.GET("/dashboard", func(c echo.Context) error {
		return dashboard.DashboardHandler(c)
	})

	authApp.POST("/api/pdf_upload", func(c echo.Context) error {
		return pdfcsv.UploadHandler(c)
	})

	authApp.GET("/api/stations", func(c echo.Context) error {
		return weights_routes.ReadWeights(c)
	})

	authApp.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics

	// 404 for routes not mounted
	authApp.GET("/*", notFound.NotFoundHandler)

	app.Logger.Fatal(app.Start(":42069"))
}
