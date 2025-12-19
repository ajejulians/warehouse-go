package app

import (
	"context"
	"warehouse-go/product-service/configs"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	

	"github.com/gofiber/fiber/v2/middleware/cors"    
	"github.com/gofiber/fiber/v2/middleware/logger"    
	"github.com/gofiber/fiber/v2/middleware/recover"

	zerolog "github.com/rs/zerolog/log"
)

func RunServer() {
	cfg := configs.NewConfig()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			zerolog.Printf("Error: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		},
	})

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New(logger.Config{  // ✅ Perbaikan sintaks
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
	
	zerolog.Info().Msgf("Starting server on port %s", port)  

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	zerolog.Info().Msg("Shutting down server...")  

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Error during shutdown: %v", err)
	}
	
	zerolog.Info().Msg("Server shutdown complete")  
}