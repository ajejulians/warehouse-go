package app

import (
	"log"
	"warehouse-go/user-service/configs"
	"warehouse-go/user-service/controller"
	"warehouse-go/user-service/database"
	"warehouse-go/user-service/pkg/storage"
	"warehouse-go/user-service/repository"
	"warehouse-go/user-service/service"
	"warehouse-go/user-service/usecase"
)

type Container struct {
	RoleController controller.RoleControllerInterface
	UserController controller.UserControllerInterface
	AuthController controller.AuthControllerInterface
	UploadController controller.UploadControllerInterface 
}

func BuildContainer() *Container {
	config := configs.NewConfig()
	db, err := database.ConnectPostgres(*config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	rabbitMQService, err := service.NewRabbitMQService(*config)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	supabaseStorage := storage.NewSupabaseStorage(*config)

	fileUploadHelper := storage.NewUploadFileHelper(supabaseStorage, *config)

	roleRepo := repository.NewRoleRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)
	roleUsecase := usecase.NewRoleUsecase(roleRepo)
	roleController := controller.NewRoleController(roleUsecase)
	UserUsecase := usecase.NewUserUsecase(userRepo, rabbitMQService)
	UserController := controller.NewUserController(UserUsecase)

	authController := controller.NewAuthController(UserUsecase)

	uploadController := controller.NewUploadController(fileUploadHelper)

	return &Container{
		RoleController: roleController,
		UserController: UserController,
		AuthController: authController,
		UploadController: uploadController,
	}
}