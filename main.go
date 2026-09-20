package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"webapp/pkg/config"
)

func main() {
	port := config.GetEnv("PORT", "3000")
	logPath := config.GetEnv("LOG_PATH", "/app/log/app.log")

	// Create the directory for the access log.
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		log.Fatalf("Create log directory: %v", err)
	}

	// Create the log file or append to the existing file.
	logFile, err := os.OpenFile(
		logPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatalf("Open log file: %v", err)
	}
	defer logFile.Close()

	e := echo.New()

	// Write HTTP request logs to the shared log file.
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: logFile,
	}))

	e.Static("/", "public")

	e.GET("/", func(c echo.Context) error {
		return c.File("public/views/webapp.html")
	})

	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Start HTTP server: %v", err)
	}
}
