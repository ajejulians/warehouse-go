package controller

import (
	"warehouse-go/merchant-service/controller/response"
	"warehouse-go/merchant-service/pkg/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UploadControllerInterface interface {
	UploadMerchantPhoto(c *fiber.Ctx) error
}

type uploadController struct {
	fileUploadHelper *storage.FileUploadHelper
}

// UploadMerchantPhoto implements UploadControllerInterface.
func (u *uploadController) UploadMerchantPhoto(c *fiber.Ctx) error {
	log.Infof("Content-Type incoming: %s", c.Get("Content-Type"))
	file, err := c.FormFile("image")
	if err != nil {
		log.Errorf("[UploadController] UploadMerchantPhoto - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message" : "No file uploaded",
			"error" : err.Error(),
		})
	}

	//upload to supabase using fileUploadHelper
	result, err := u.fileUploadHelper.UploadPhoto(c.Context(), file)
	if err != nil {
		log.Errorf("[UploadController] UploadMerchantPhoto - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" : "Failed to upload file",
			"error" : err.Error(),
		})
	}

	//Create response
	uploadResponse := response.UploadResponse{
		URL: result.URL,
		Filename: result.Filename,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message" : "File uploaded successfully",
		"data" : uploadResponse,
	})
}

func NewUploadController(fileUploadHelper *storage.FileUploadHelper) UploadControllerInterface {
	return &uploadController{
		fileUploadHelper: fileUploadHelper,
	}
}
