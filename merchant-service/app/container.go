package app

import (
	"log"
	"warehouse-go/merchant-service/configs"
	"warehouse-go/merchant-service/controller"
	"warehouse-go/merchant-service/database"
	"warehouse-go/merchant-service/pkg/httpclient"
	"warehouse-go/merchant-service/pkg/rabbitmq"
	"warehouse-go/merchant-service/pkg/redis"
	"warehouse-go/merchant-service/pkg/storage"
	"warehouse-go/merchant-service/repository"
	"warehouse-go/merchant-service/usecase"
)

type Container struct {
	MerchantController controller.MerchantControllerInterface
	MerchantProductController controller.MerchantProductControllerInterface
	UploadController controller.UploadControllerInterface
}

func BuildContainer() *Container {
	cfg := configs.NewConfig()
	db, err := database.ConnectPostgres(*cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	redisClient := redis.NewRedisClient(*cfg)
	rabbitMQService, err := rabbitmq.NewRabbitMQService(cfg.RabbitMQ.URL())
	if err != nil {
		log.Fatalf("Failed to connect to rabbitmq: %v", err)	
	}

	userClient := httpclient.NewUserClient(*cfg)
	cachedUserClient := httpclient.NewCachedUserClient(userClient, redisClient)
	warehouseClient := httpclient.NewWarehouseClient(*cfg)
	cachedWarehouseClient := httpclient.NewCachedWarehouseClient(warehouseClient, redisClient)
	productClient := httpclient.NewProductClient(*cfg)
	cachedProductClient := httpclient.NewCachedProductClient(productClient, redisClient)

	merchantRepo := repository.NewMerchantRepository(db.DB)
	merchantUsecase := usecase.NewMerchantUsecase(merchantRepo, cachedUserClient, cachedWarehouseClient, cachedProductClient)
	merchantController := controller.NewMerchantController(merchantUsecase)

	merchantProductRepo := repository.NewMerchantProductRepository(db.DB)
	merchantProductUsecase := usecase.NewMerchantProductUsecase(merchantProductRepo, cachedProductClient, cachedWarehouseClient, rabbitMQService)
	merchantProductController := controller.NewMerchantProductController(merchantProductUsecase)

	supabaseStorage := storage.NewSupabaseStorage(*cfg)
	uploadFileHelper := storage.NewUploadFileHelper(supabaseStorage, *cfg)
	uploadController := controller.NewUploadController(uploadFileHelper)	

	return &Container {
		MerchantController: merchantController,
		MerchantProductController: merchantProductController,
		UploadController: uploadController,
	}
}