package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	jwtConf "warehouse-go/api-gateaway/config"
	"warehouse-go/api-gateaway/controller"
	"warehouse-go/api-gateaway/middleware"
)


type ServiceConfig struct {
	Name string
	URL string
}

type Config struct {
	Port     string
	Services map[string]ServiceConfig
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found, trying env.example")
		if err := godotenv.Load("env.example"); err != nil {
			log.Println("Warning: env.example file not found, using default configuration")
		}
	}

	config := loadConfig()
	jwtConfig := jwtConf.LoadJWTConfig()
	redisConfig := jwtConf.LoadRedisConfig()

	app := fiber.New(fiber.Config{
		AppName: "Warehouse Project API Gateaway",
		ServerHeader: "Warehouse-API-Gateaway",
	})

	rateLimiter := middleware.DefaultRateLimiterConfig()

	var redisClient *redis.Client
	if redisConfig.Host != "" {
		redisClient = jwtConf.NewRedisClient(redisConfig)
	}
	redisRateConfig := middleware.RedisRateLimiterConfig{
		Max: rateLimiter.Max,
		Expiration: rateLimiter.Expiration,
		RedisClient: redisClient,
	}

	app.Use(middleware.RedisGlobalRateLimiter(redisRateConfig))
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Auhtorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Get("/health", func (c *fiber.Ctx) error{
		return c.JSON(fiber.Map{
			"status" : "OK",
			"message" : "API Gateaway is running",
			"services" : fiber.Map{
				"user-services" : config.Services["user"].URL,
				"product-services" : config.Services["product"].URL,
				"merchant-services" : config.Services["merchant"].URL,
				"notification-services" : config.Services["notification"].URL,
				"transaction-services" : config.Services["transaction"].URL,
				"midtrans-services" : config.Services["midtrans"].URL,
				"warehouse-services" : config.Services["warehouse"].URL,
				"role-services" : config.Services["role"].URL,
				"assign-role-services" : config.Services["assign-role"].URL,
			},
		})
	})

	authController := controller.NewAuthController(config.Services["auth"].URL, jwtConfig)
	setUpAuthRoutes(app, authController, rateLimiter)
	setUpMidtransCallbackRoutes(app, config.Services["midtrans"])

	setupProtectedRoutes(app, config, jwtConfig, rateLimiter)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"error" : "Not Found",
			"message" : "Service Not Found",
			"path" : c.Path(),
		})
	})

	log.Printf("API Gateaway starting on port %s", config.Port)
	log.Fatal(app.Listen(":" + config.Port))
}

func loadConfig() Config {
	config := Config{
		Port: getEnv("PORT", "8080"),
		Services: map[string]ServiceConfig{
			"user" : {
				Name: "user-service",
				URL: getEnv("USER_SERVICE_URL", "http://localhost:8081"),
			},
			"role" : {
				Name: "role-service",
				URL: getEnv("USER_SERVICE_URL", "http://localhost:8081"),
			},
			"assign-role" : {
				Name: "assign-role-service", 
				URL: getEnv("USER_SERVICE_URL", "http://localhost:8081"),
			},
			"auth" : {
				Name: "auth-service", 
				URL: getEnv("USER_SERVICE_URL", "http://localhost:8081"),
			},
			"product" : {
				Name: "product-service", 
				URL: getEnv("PRODUCT_SERVICE_URL", "http://localhost:8082"),
			},
			"merchant" : {
				Name: "merchant-service", 
				URL: getEnv("MERCHANT_SERVICE_URL", "http://localhost:8084"),
			},
			"notification" : {
				Name: "notification-service", 
				URL: getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8086"),
			},
			"transaction" : {
				Name: "transaction-service", 
				URL: getEnv("TRANSACTION_SERVICE_URL", "http://localhost:8085"),
			},
			"midtrans" : {
				Name: "midtrans-service", 
				URL: getEnv("TRANSACTION_SERVICE_URL", "http://localhost:8085"),
			},
			"warehouse" : {
				Name: "warehouse-service", 
				URL: getEnv("WAREHOUSE_SERVICE_URL", "http://localhost:8083"),
			},
		},
	}

	log.Println("Service Configuration")
	for name, service := range config.Services {
		log.Printf(" %s: %s", name, service.URL)
	}

	return config
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func setUpAuthRoutes(app *fiber.App, authController *controller.AuthController, rateLimiterConfig middleware.RedisRateLimiterConfig) {
	authGroup := app.Group("/api/v1/auth")

	authGroup.Use(middleware.RedisAuthRateLimiter(rateLimiterConfig))
	authGroup.Post("/login", authController.Login)
}

func setUpMidtransCallbackRoutes(app *fiber.App, service ServiceConfig) {
	app.Post("/api/v1/midtrans/callback", func(c *fiber.Ctx) error {
		client := &http.Client{}

		fullUrl := service.URL + "/api/v1/midtrans/callback"

		body := c.Body()

		req, err := http.NewRequest(c.Method(), fullUrl, bytes.NewReader(body))
		if err != nil {
			log.Printf("Error createing request: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error" : "Internal Server Error",
				"message" : "Failed to create request",
			})
		}

		for key, values := range c.GetReqHeaders() {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		req.Header.Set("X-Gateaway", "warehouse-api-gateaway")
		req.Header.Set("X-Internal-Request", "true")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error making request to: %s: %v", fullUrl, err)
			return c.Status(502).JSON(fiber.Map{
				"error" : "Bad Request",
				"message" : "Service unavailable",
				"service" : service.URL,
			})
		}

		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error" : "Internal Server Error",
				"message" : "Failed to read response",
			})
		}
		
		for key, values := range resp.Header {
			for _, value := range values {
				c.Set(key, value)
			}
		}

		return c.Status(resp.StatusCode).Send(respBody)

	})
}

func proxyRequestWithPath(c *fiber.Ctx, targetURL string, basePath string) error {
	fullPath := c.Path()
	
	fullURL := targetURL + fullPath

	queryParams := c.Context().QueryArgs().String()
	if queryParams != "" {
		fullURL += "?" + queryParams
	}

	client := &http.Client{}

	body := c.Body()

	req, err := http.NewRequest(c.Method(), fullURL, bytes.NewReader(body))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error" : "Internal Server Error",
			"message" : "Failed to create request",
		})
	}

	for key, values := range c.GetReqHeaders() {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	req.Header.Set("X-Gateaway", "warehouse-api-gateaway")
	req.Header.Set("X-Internal-Request", "true")
	
	addUserHeaders(req, c) 

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to: %s: %v", fullURL, err )
		return c.Status(502).JSON(fiber.Map{
			"error" : "Bad Gateaway",
			"message" : "Service Unavailable",
			"service" : targetURL, 
		})
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error" : "Internal Server Error",
			"message" : "Failed to read response",
		})
	}

	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	return c.Status(resp.StatusCode).Send(respBody)
}

func proxyRequest(c *fiber.Ctx, targetURL string) error {
	path := c.Params("*")
	if path == "" {
		path = c.Path()
	}

	path = strings.TrimPrefix(path, "api/v1/")

	fullURL := targetURL
	if !strings.HasSuffix(targetURL, "/") && !strings.HasPrefix(path, "/") {
		fullURL += "/"
	}
	fullURL += path

	queryParams := c.Context().QueryArgs().String()
	if queryParams != "" {
		fullURL += "?" + queryParams
	}

	client := &http.Client{}

	body := c.Body()

	req, err := http.NewRequest(c.Method(), fullURL, bytes.NewReader(body))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error" : "Internal Server Error",
			"message" : "Failed to create request",
		})
	}

	for key, values := range c.GetReqHeaders() {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	req.Header.Set("X-Gateaway", "warehouse-api-gateaway")
	req.Header.Set("X-Internal-Request", "true")
	
	addUserHeaders(req, c) 

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to: %s: %v", fullURL, err )
		return c.Status(502).JSON(fiber.Map{
			"error" : "Bad Gateaway",
			"message" : "Service Unavailable",
			"service" : targetURL, 
		})
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error" : "Internal Server Error",
			"message" : "Failed to read response",
		})
	}

	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	return c.Status(resp.StatusCode).Send(respBody)

}

func addUserHeaders( req *http.Request, c *fiber.Ctx) {
	if userID := c.Locals("user_id"); userID != nil {
		req.Header.Set("X-User-ID", fmt.Sprintf("%v", userID))
	}
	if userEmail := c.Locals("user_email"); userEmail != nil {
		req.Header.Set("X-User-email", fmt.Sprintf("%v", userEmail))
	}
	if userRoles := c.Locals("user_roles"); userRoles != nil {
		req.Header.Set("X-User-Roles", fmt.Sprintf("%v", userRoles))
	}
}

func setupProtectedRoutes(app *fiber.App, config Config, jwtConfig middleware.JWTConfig, rateLimiterConfig middleware.RedisRateLimiterConfig) {
	protected := app.Group("/api/v1", middleware.JWTAuthMiddleware(jwtConfig))

	protected.Use(middleware.RedisAPIRateLimiter(rateLimiterConfig))

	setupUserRoutes(protected, config.Services["user"])
	setupRolesRoutes(protected, config.Services["role"])
	setupAssignRoleRoutes(protected, config.Services["assign-role"])
	setupProductRoutes(protected, config.Services["product"])
	setupMerchantRoutes(protected, config.Services["merchant"])
	setupTransactionRoutes(protected, config.Services["transaction"])
	setupWarehouseRoutes(protected, config.Services["warehouse"])


}

func setupUserRoutes(router fiber.Router, service ServiceConfig) {
	userGroup := router.Group("/users")

	userGroup.All("/*", func(c *fiber.Ctx) error {
		 return proxyRequestWithPath(c, service.URL, "/api/v1/users")
	})

	userGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/users")
	})
}

func setupRolesRoutes(router fiber.Router, service ServiceConfig) {
	roleGroup := router.Group("/roles")

	roleGroup.All("/*", func(c *fiber.Ctx) error {
		 return proxyRequestWithPath(c, service.URL, "/api/v1/roles")
	})

	roleGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/roles")
	})
}

func setupAssignRoleRoutes(router fiber.Router, service ServiceConfig) {
	assignRoleGroup := router.Group("/assign-role")

	assignRoleGroup.All("/*", func(c *fiber.Ctx) error {
		 return proxyRequestWithPath(c, service.URL, "/api/v1/assign-role")
	})

	assignRoleGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/assign-role")
	})
}

func setupProductRoutes(router fiber.Router, service ServiceConfig) {
	productGroup := router.Group("/product")

	productGroup.All("/*", func(c *fiber.Ctx) error {
		 return proxyRequestWithPath(c, service.URL, "/api/v1/products")
	})

	productGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/products")
	})

	categoryGroup := router.Group("/categories")

	categoryGroup.All("/*", func(c *fiber.Ctx) error {
		 return proxyRequestWithPath(c, service.URL, "/api/v1/categories")
	})

	categoryGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/categories")
	})

	uploadProductGroup := router.Group("/upload")

	uploadProductGroup.All("/*", func(c *fiber.Ctx) error {
		 return proxyRequestWithPath(c, service.URL, "/api/v1/upload")
	})

	uploadProductGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/produploaducts")
	})
}

func setupMerchantRoutes(router fiber.Router, service ServiceConfig) {
	merchantGroup := router.Group("/merchants")

	merchantGroup.All("/*", func(c *fiber.Ctx) error {
		 return proxyRequestWithPath(c, service.URL, "/api/v1/merchants")
	})

	merchantGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/merchants")
	})

	merchantProductGroup := router.Group("/merchant-products")

	merchantProductGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/merchant-products")
	})

	merchantProductGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/merchant-products")
	})

	uploadGroup := router.Group("/upload-merchant")
	uploadGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequest(c, service.URL)
	})

	uploadGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequest(c, service.URL)
	})
} 

func setupTransactionRoutes(router fiber.Router, service ServiceConfig) {
	transactionGroup := router.Group("/transactions")

	transactionGroup.All("/*", func(c *fiber.Ctx) error {
		 return proxyRequestWithPath(c, service.URL, "/api/v1/transactions")
	})

	transactionGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/transactions")
	})

	dashboardGroup := router.Group("/dashboard")

	dashboardGroup.All("/*", func(c *fiber.Ctx) error {
		 return proxyRequestWithPath(c, service.URL, "/api/v1/dashboard")
	})

	dashboardGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/dashboard")
	})
}

func setupWarehouseRoutes(router fiber.Router, service ServiceConfig) {
	warehouseGroup := router.Group("/warehouses")

	warehouseGroup.All("/*", func(c *fiber.Ctx) error {
		 return proxyRequestWithPath(c, service.URL, "/api/v1/warehouses")
	})

	warehouseGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/warehouses")
	})

	warehouseProductGroup := router.Group("/warehouse-products")

	warehouseProductGroup.All("/*", func(c *fiber.Ctx) error {
		 return proxyRequestWithPath(c, service.URL, "/api/v1/warehouse-products")
	})

	warehouseProductGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/warehouse-products")
	})

	uploadWarehouseGroup := router.Group("/upload-warehouse")

	uploadWarehouseGroup.All("/*", func(c *fiber.Ctx) error {
		 return proxyRequest(c, service.URL)
	})

	uploadWarehouseGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequest(c, service.URL)
	})
}




