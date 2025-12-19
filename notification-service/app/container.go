package app

import (
	"warehouse-go/notification-service/controller"
	"warehouse-go/notification-service/pkg/email"
	"warehouse-go/notification-service/pkg/rabbitmq"
	"warehouse-go/notification-service/usecase"
)

type Container struct {
	EmailController controller.EmailController
	EmailUseCase usecase.EmailUseCase
	RabbitMQService rabbitmq.RabbitMQServiceInterface
	EmailService email.EmailServiceInterface
}

func BuildContainer(rabbitMQService rabbitmq.RabbitMQServiceInterface, emailService email.EmailServiceInterface) *Container {
	emailUseCase := usecase.NewEmailUsecase(emailService)
	emailController := controller.NewEmailController(emailUseCase)

	return &Container{
		EmailController: *emailController,
		EmailUseCase: *emailUseCase,
		RabbitMQService: rabbitMQService,
		EmailService: emailService,
	}
}