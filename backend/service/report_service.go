package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/dto"
	"github.com/Caknoooo/go-gin-clean-starter/entity"
	"github.com/Caknoooo/go-gin-clean-starter/repository"
	"github.com/Caknoooo/go-gin-clean-starter/utils"
)

type (
	ReportService interface {
		CreateReport(ctx context.Context, req dto.CreateReportRequest) (dto.CreateReportResponse, error)
	}

	reportService struct {
		userRepo   repository.UserRepository
		reportRepo repository.ReportRepository
		db         *gorm.DB
	}
)

func NewReportService(
	userRepo repository.UserRepository,
	reportRepo repository.ReportRepository,
	db *gorm.DB,
) ReportService {
	return &reportService{
		userRepo:   userRepo,
		reportRepo: reportRepo,
		db:         db,
	}
}

// Constants for image validation
const (
	MaxImageSize   = 10 * 1024 * 1024 // 10MB
	MaxImageWidth  = 4000             // pixels
	MaxImageHeight = 4000             // pixels
	UploadDir      = "./uploads/reports"
)

// validateImageSize checks if image size is within allowed limits
func validateImageSize(file *multipart.FileHeader) error {
	if file.Size > MaxImageSize {
		return dto.ErrImageSizeTooLarge
	}

	if file.Size == 0 {
		return dto.ErrEmptyFileUploaded
	}

	return nil
}

// isValidImageType validates image type by checking both MIME type and file signature
func isValidImageType(file *multipart.FileHeader) (bool, error) {
	// Define allowed MIME types
	allowedMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
		"image/bmp":  true,
		"image/tiff": true,
	}

	// Check MIME type from header
	contentType := file.Header.Get("Content-Type")
	if !allowedMimeTypes[strings.ToLower(contentType)] {
		return false, dto.ErrInvalidImageType
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".bmp":  true,
		".tiff": true,
		".tif":  true,
	}

	if !allowedExtensions[ext] {
		return false, dto.ErrInvalidImageExtension
	}

	// Verify file signature (magic bytes) for additional security
	src, err := file.Open()
	if err != nil {
		return false, dto.ErrFailedToOpenImageFile
	}
	defer src.Close()

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return false, dto.ErrFailedToOpenImageFile
	}

	// Reset file pointer
	src.Seek(0, 0)

	// Detect actual content type from file content
	detectedType := http.DetectContentType(buffer)

	// Verify detected type matches allowed types
	if !allowedMimeTypes[strings.ToLower(detectedType)] {
		return false, dto.ErrInvalidImageType
	}

	return true, nil
}

func (s *reportService) CreateReport(ctx context.Context, req dto.CreateReportRequest) (dto.CreateReportResponse, error) {
	user_id, ok := ctx.Value("user_id").(string)
	if !ok {
		return dto.CreateReportResponse{}, dto.ErrUserIdEmpty
	}

	if req.Image == nil && req.Text == "" {
		return dto.CreateReportResponse{}, dto.ErrEmptyContent
	}

	// Validate user exists
	_, err := s.userRepo.GetUserById(ctx, nil, user_id)
	if err != nil {
		return dto.CreateReportResponse{}, dto.ErrUserNotFound
	}

	reportID := uuid.New().String()
	var imagePath string

	// Handle Image validation and processing
	if req.Image != nil {
		ext := utils.GetExtensions(req.Image.Filename)
		imagePath = fmt.Sprintf("reports/%s.%s", reportID, ext)
		if err := utils.UploadFile(req.Image, imagePath); err != nil {
			return dto.CreateReportResponse{}, err
		}
	}

	// Create report entity (prepared for database insertion)
	report := entity.Report{
		ID:     reportID,
		Text:   req.Text,
		Image:  imagePath,
		UserID: user_id,
		Status: entity.StatusUnverified,
		TagID:  nil,
	}

	createdReport, err := s.reportRepo.CreateReport(ctx, nil, report)
	if err != nil {
		if imagePath != "" {
			os.Remove(imagePath)
		}
		return dto.CreateReportResponse{}, dto.ErrCreateReport
	}

	return dto.CreateReportResponse{
		ID:    createdReport.ID,
		Text:  createdReport.Text,
		Image: createdReport.Image,
	}, nil
}
