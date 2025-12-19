package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"warehouse-go/merchant-service/configs"

	"github.com/gofiber/fiber/v2/log"
	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseInterface interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*UploadResult, error)
}

type SupabaseStorage struct {
	client *storage_go.Client
	cfg    configs.Config
}

// UploadFile implements SupabaseInterface.
func (s *SupabaseStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*UploadResult, error) {
	// Log config untuk debugging
	log.Infof("Upload Config - URL: %s, Bucket: %s", s.cfg.Supabase.Url, s.cfg.Supabase.Bucket)

	// Open file
	src, err := file.Open()
	if err != nil {
		log.Errorf("Failed to open file: %v", err)
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Read file content into bytes
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		log.Errorf("Failed to read file: %v", err)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	log.Infof("File read successfully, size: %d bytes", len(fileBytes))

	// Generate file path
	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%s_%d%s", strings.TrimSuffix(file.Filename, ext), timestamp, ext)
	filePath := fmt.Sprintf("%s/%s", folder, filename)

	log.Infof("Generated file path: %s", filePath)

	// Determine content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".webp":
			contentType = "image/webp"
		case ".svg":
			contentType = "image/svg+xml"
		default:
			contentType = "application/octet-stream"
		}
	}

	log.Infof("Content-Type: %s", contentType)

	// Upload File with proper options
	upsert := true
	fileReader := bytes.NewReader(fileBytes)
	
	response, err := s.client.UploadFile(s.cfg.Supabase.Bucket, filePath, fileReader, storage_go.FileOptions{
		ContentType: &contentType,
		Upsert:      &upsert,
	})

	if err != nil {
		log.Errorf("Upload error: %v (type: %T)", err, err)
		return nil, fmt.Errorf("failed to upload to supabase: %w", err)
	}

	log.Infof("Upload successful! Response: %+v", response)

	// Get public URL
	publicUrl := s.client.GetPublicUrl(s.cfg.Supabase.Bucket, filePath)

	log.Infof("Public URL generated: %s", publicUrl.SignedURL)

	return &UploadResult{
		URL:      publicUrl.SignedURL,
		Path:     filePath,
		Filename: filename,
	}, nil
}

type UploadResult struct {
	URL      string `json:"url"`
	Path     string `json:"path"`
	Filename string `json:"filename"`
}

func NewSupabaseStorage(cfg configs.Config) SupabaseInterface {
	log.Infof("Initializing Supabase Storage - URL: %s, Bucket: %s", cfg.Supabase.Url, cfg.Supabase.Bucket)
	
	client := storage_go.NewClient(cfg.Supabase.Url, cfg.Supabase.Key, nil)
	
	return &SupabaseStorage{
		client: client,
		cfg:    cfg,
	}
}