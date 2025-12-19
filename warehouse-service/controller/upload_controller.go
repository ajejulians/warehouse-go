package controller

import (
	"warehouse-go/warehouse-service/controller/response"
	"warehouse-go/warehouse-service/pkg/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UploadControllerInterface interface {
	UploadPhoto(c *fiber.Ctx) error
}

type uploadController struct {
	fileUploadHelper *storage.FileUploadHelper
}

// UploadPhoto implements UploadControllerInterface.
func (u *uploadController) UploadPhoto(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		log.Errorf("[UploadController] UploadPhoto - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message" : "Failed to get file",
		})
	}

	result, err := u.fileUploadHelper.UploadPhoto(c.Context(), file)
	if err != nil {
		log.Errorf("[UploadController] UploadPhoto - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" : "Failed to upload file",
		})		
	}

	resp := response.UploadResponse {
		URL: result.URL,
		Path: result.Path,
		Filename: result.Filename,
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message" : "File uploaded successfully",
		"data" : resp, 
	})
}

func NewFileUploadController(fileUploadHelper *storage.FileUploadHelper) UploadControllerInterface {
	return &uploadController{fileUploadHelper: fileUploadHelper}
}
