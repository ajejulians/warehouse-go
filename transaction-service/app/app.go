package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"warehouse-go/transaction-service/configs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"    
	"github.com/gofiber/fiber/v2/middleware/logger"    
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
)

func RunServer() {
	cfg := configs.NewConfig()

	zlog := zerolog.New(os.Stderr).With().Timestamp().Logger()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			zlog.Error().Err(err).Msg("Error")
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		},
	})

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New(logger.Config{  
		Format: "[${time}] ${ip} ${status} - ${method} ${path}\n",  
	}))

	container := BuildContainer()
	SetupRoutes(app, container)

	port := cfg.App.AppPort
	if port == "" {
		port = os.Getenv("APP_PORT")
		if port == "" {
			log.Fatal("Server port not specified")
		}
	}
	
	zlog.Info().Msgf("Starting server on port %s", port)  

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	zlog.Info().Msg("Shutting down server...")  

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Error during shutdown: %v", err)
	}
	
	zlog.Info().Msg("Server shutdown complete")  
}