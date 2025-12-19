package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"warehouse-go/product-service/configs"


	"github.com/gofiber/fiber/v2/log"
)

const (
	MaxImageSize = 2 * 1024 * 1024 // 2MB

	AllowedImageExtensions = ".jpg, .jpeg, .png, .webp, .svg"
)

type FileUploadHelper struct {
	storage SupabaseInterface
	cfg     configs.Config
}

func NewUploadFileHelper(storage SupabaseInterface, cfg configs.Config) *FileUploadHelper {
	return &FileUploadHelper{
		storage: storage, 
		cfg: cfg,
	}
}

func (h *FileUploadHelper) UploadPhoto(ctx context.Context, file *multipart.FileHeader, folder string) (*UploadResult, error){
	if err := h.validateImageFile(file, MaxImageSize); err != nil {
		log.Errorf("failed to validate image file: %v", err)
		return nil, err
	}

	result, err := h.storage.UploadFile(ctx, file, folder)
	if err != nil{
		log.Errorf("failed to upload file: %v", err)
		return nil, err
	}

	return  result, nil

}

func (h *FileUploadHelper) validateImageFile(file *multipart.FileHeader, maxSize int64) error{
	if !validateFileSize(file.Size, maxSize) {
		return fmt.Errorf("file size exceeds the maximum allowed size")
	}

	if !validateFileExtension(getFileExtension(file.Filename), AllowedImageExtensions) {
		return fmt.Errorf("invalid file extension")
	}

	return nil
}

func validateFileSize(size int64, maxSize int64) bool {
	return size <= maxSize
}

func getFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

func validateFileExtension(extension string, allowedExtensions string) bool{
	allowed := strings.Split(allowedExtensions, ",")
	for _, ext := range allowed {
		if strings.TrimSpace(ext) == extension {
			return true
		}
	}
	return false
}