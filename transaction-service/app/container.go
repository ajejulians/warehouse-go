package app

import (
	"log"
	"warehouse-go/transaction-service/configs"
	"warehouse-go/transaction-service/controller"
	"warehouse-go/transaction-service/database"
	"warehouse-go/transaction-service/pkg/httpclient"
	"warehouse-go/transaction-service/pkg/midtrans"
	"warehouse-go/transaction-service/pkg/rabbitmq"
	"warehouse-go/transaction-service/repository"
	"warehouse-go/transaction-service/usecase"
)

type Container struct {
	TransactionController controller.TransactionControllerInterface
}

func BuildContainer() *Container {
	cfg := configs.NewConfig()

	db, err := database.ConnectPostgres(*cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	transactionRepo := repository.NewTransactionRepository(db.DB)

	//HTTP Clients
	merchantClient := httpclient.NewMerchantClient(*cfg)
	userClient := httpclient.NewUserClient(*cfg)
	productClient := httpclient.NewProductClient(*cfg)

	//RabbitMQ Client
	rabbitMQService, err := rabbitmq.NewRabbitMQService(cfg.RabbitMQ.URL())
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo, merchantClient, rabbitMQService, productClient, userClient )
	midtransService := midtrans.NewMidtransService(cfg)
	transactionController := controller.NewTransactionController(transactionUsecase, midtransService)
	
	return &Container{
		TransactionController: transactionController,
	}
}