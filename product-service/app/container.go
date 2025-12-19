package app

import (
	"warehouse-go/product-service/configs"
	"warehouse-go/product-service/controller"
	"warehouse-go/product-service/database"
	"warehouse-go/product-service/pkg/storage"
	"warehouse-go/product-service/repository"
	"warehouse-go/product-service/usecase"

	"github.com/gofiber/fiber/v2/log"
)

type Container struct {
	ProductController controller.ProductControllerInterface
	CategoryController controller.CategoryControllerInterface
	UploadController controller.UploadControllerInterface

}

func BuildContainer() *Container {
	config := configs.NewConfig()
	db, err := database.ConnectPostgres(*config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	categoryRepo := repository.NewCategoryRepository(db.DB)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	categoryController := controller.NewCategoryController(categoryUsecase)

	productRepo := repository.NewProductRepository(db.DB)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productController := controller.NewProductController(productUsecase)

	supabaseStorage := storage.NewSupabaseStorage(*config)
	fileUploadHelper := storage.NewUploadFileHelper(supabaseStorage, *config)
	uploadController := controller.NewUploadController(fileUploadHelper)

	return &Container{
		ProductController: productController,
		CategoryController: categoryController,
		UploadController: uploadController,
	}
}
