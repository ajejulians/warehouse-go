package app

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App, container *Container) {
	api := app.Group("/api/v1")
	categories := api.Group("/categories")
	products := api.Group("/products")
	uploads := api.Group("/upload")

	categories.Post("/", container.CategoryController.CreateCategory)
	categories.Get("/", container.CategoryController.GetAllCategories)
	categories.Get("/:id", container.CategoryController.GetCategoryByID)
	categories.Put("/:id", container.CategoryController.UpdateCategory)
	categories.Delete("/:id", container.CategoryController.DeleteCategory)


	products.Post("/", container.ProductController.CreateProduct)
	products.Get("/", container.ProductController.GetAllProducts)
	products.Get("/:id", container.ProductController.GetProductByID)
	products.Get("/barcode/:barcode", container.ProductController.GetProductByBarcode)
	products.Put("/:id", container.ProductController.UpdateProduct)
	products.Delete("/:id", container.ProductController.DeleteProduct)

	uploads.Post("/product", container.UploadController.UploadProductImage )
	uploads.Post("/category-image", container.UploadController.UploadCategoryImage)
}