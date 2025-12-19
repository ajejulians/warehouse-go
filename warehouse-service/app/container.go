package app

import (
	"log"
	"time"
	"warehouse-go/warehouse-service/configs"
	"warehouse-go/warehouse-service/controller"
	"warehouse-go/warehouse-service/database"
	"warehouse-go/warehouse-service/pkg/httpclient"
	"warehouse-go/warehouse-service/pkg/rabbitmq"
	"warehouse-go/warehouse-service/pkg/redis"
	"warehouse-go/warehouse-service/pkg/storage"
	"warehouse-go/warehouse-service/repository"
	"warehouse-go/warehouse-service/usecase"
)

type Container struct {
	WarehouseController controller.WarehouseControllerInterface
	WarehouseProductController controller.WarehouseProductControllerInterface
	UploadController controller.UploadControllerInterface
	RabbitMQConsumer *rabbitmq.RabbitMQConsumer
}

func BuildContainer() *Container {
	config := configs.NewConfig()
	db, err := database.ConnectPostgres(*config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	productClient := httpclient.NewProductClient(*config)
	redisClient := redis.NewRedisClient(*config)
	cacheProductClient := httpclient.NewCachedProductClient(productClient, redisClient, 1 *time.Hour)

	warehouseRepo := repository.NewWarehouseRepository(db.DB)
	warehouseUsecase := usecase.NewWarehouseUsecase(warehouseRepo)
	warehouseController := controller.NewWarehouseController(warehouseUsecase)

	warehouseProductRepo := repository.NewWarehouseProductRepository(db.DB)
	warehouseProductUsecase := usecase.NewWarehouseProductUsecase(warehouseProductRepo, cacheProductClient)
	warehouseProductController := controller.NewWarehouseProductController(warehouseProductUsecase)

	rabbitMQConsumer, err := rabbitmq.NewRabbitMQConsumer(config.RabbitMQ.URL(), warehouseProductRepo)
	if err != nil {
		log.Fatalf("Failed to create rabbitmq consumer: %v", err)
	}


	supabaseStorage := storage.NewSupabaseStorage(*config)
	fileUploadHelper := storage.NewUploadFileHelper(supabaseStorage, *config)
	uploadController := controller.NewFileUploadController(fileUploadHelper)

	 

	return &Container{
		WarehouseController: warehouseController,
		WarehouseProductController: warehouseProductController,
		UploadController: uploadController,
		RabbitMQConsumer: rabbitMQConsumer,
	}
}