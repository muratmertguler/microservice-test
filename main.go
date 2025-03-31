package main

import (
	"context"
	"fmt"
	"microservicetest/app/healthcheck"
	"microservicetest/pkg/config"
	_ "microservicetest/pkg/log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Request any
type Response any

type HandlerInterface[R Request, Res Response] interface {
	Handle(ctx context.Context, Req *R) (*Res, error)
}

func handle[R Request, Res Response](handler HandlerInterface[R, Res]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req R

		if c.Method() != fiber.MethodGet {
			if err := c.BodyParser(&req); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
			}
		}

		res, err := handler.Handle(c.Context(), &req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
		}
		return c.JSON(res)
	}
}

func main() {
	appConfig := config.Read()
	fmt.Printf("appConfig: %+v\n", appConfig)
	defer zap.L().Sync()

	zap.L().Info("starting server...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://www.google.com", nil)
	if err != nil {
		zap.L().Error("failed to create request", zap.Error(err))
	}
	healthcheckHandler := healthcheck.NewHealthCheckHandler()

	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		Concurrency:  256 * 1024,
	})

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	app.Get("/", func(c *fiber.Ctx) error {
		zap.L().Info("server started")
		return c.SendString("Hello World")
	})
	app.Get("/healthcheck", handle[healthcheck.HealthCheckRequest, healthcheck.HealthCheckResponse](healthcheckHandler))

	// Starting server in a goroutine
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", appConfig.Port)); err != nil {
			zap.L().Error("failed to start server", zap.Error(err))
			os.Exit(1)
		}
	}()

	zap.L().Info("server started", zap.String("port", appConfig.Port))

	// Graceful-Shutdown-Server
	gracefulShutdown(app)
}

func gracefulShutdown(app *fiber.App) {
	// Create channel for shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, os.Interrupt)

	// Wait for shutdown signal
	<-sigChan
	zap.L().Info("shutting down server...")

	app.Shutdown()
	if err := app.ShutdownWithTimeout(2 * time.Second); err != nil {
		zap.L().Error("Error during shutdown server", zap.Error(err))
	}

	zap.L().Info("server gracefully stopped")
}

func httpc() {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		zap.L().Error("failed to create request", zap.Error(err))
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		zap.L().Error("failed to do request", zap.Error(err))
	}

	zap.L().Info("google response ", zap.Int("status", resp.StatusCode))

}
