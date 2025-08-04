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
		GetAllReports(ctx context.Context, req dto.PaginationRequest) (dto.ReportPaginationResponse, error)
		GetReportById(ctx context.Context, reportId string) (dto.ReportResponse, error)
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
		ID:       reportID,
		Text:     req.Text,
		Image:    imagePath,
		UserID:   user_id,
		Status:   entity.StatusUnverified,
		Location: req.Location,
		TagID:    nil,
	}

	createdReport, err := s.reportRepo.CreateReport(ctx, nil, report)
	if err != nil {
		if imagePath != "" {
			os.Remove(imagePath)
		}
		return dto.CreateReportResponse{}, dto.ErrCreateReport
	}

	return dto.CreateReportResponse{
		ID:       createdReport.ID,
		Text:     createdReport.Text,
		Image:    createdReport.Image,
		Location: createdReport.Location,
	}, nil
}

func (s *reportService) GetAllReports(ctx context.Context, req dto.PaginationRequest) (dto.ReportPaginationResponse, error) {
	reports, err := s.reportRepo.GetAllReportsWithPagination(ctx, nil, req)
	if err != nil {
		return dto.ReportPaginationResponse{}, dto.ErrGetReports
	}

	var datas []dto.ReportResponse
	for _, report := range reports.Reports {
		user, err := s.userRepo.GetUserById(ctx, nil, report.UserID)
		if err != nil {
			return dto.ReportPaginationResponse{}, dto.ErrGetReportById
		}
		data := dto.ReportResponse{
			ID:         report.ID,
			Text:       report.Text,
			Image:      report.Image,
			Location:   report.Location,
			Status:     fmt.Sprintf("%v", report.Status),
			Upvotes:    report.Upvotes,
			ShareCount: report.ShareCount,
			UserID:     report.UserID,
			Username:   user.Name,
			TagID: func() string {
				if report.TagID != nil {
					return *report.TagID
				}
				return ""
			}(),
			PredConfidence: func() int {
				if report.PredConfidence != nil {
					return *report.PredConfidence
				}
				return 0
			}(),
		}
		datas = append(datas, data)
	}

	return dto.ReportPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    reports.Page,
			PerPage: reports.PerPage,
			MaxPage: reports.MaxPage,
			Count:   reports.Count,
		},
	}, nil
}

func (s *reportService) GetReportById(ctx context.Context, reportId string) (dto.ReportResponse, error) {
	report, err := s.reportRepo.GetReportById(ctx, nil, reportId)
	if err != nil {
		return dto.ReportResponse{}, dto.ErrGetReportById
	}
	user, err := s.userRepo.GetUserById(ctx, nil, report.UserID)
	if err != nil {
		return dto.ReportResponse{}, dto.ErrGetUserById
	}

	return dto.ReportResponse{
		ID:         report.ID,
		Text:       report.Text,
		Image:      report.Image,
		Location:   report.Location,
		Status:     fmt.Sprintf("%v", report.Status),
		Upvotes:    report.Upvotes,
		ShareCount: report.ShareCount,
		UserID:     report.UserID,
		Username:   user.Name,
		TagID: func() string {
			if report.TagID != nil {
				return *report.TagID
			}
			return ""
		}(),
		PredConfidence: func() int {
			if report.PredConfidence != nil {
				return *report.PredConfidence
			}
			return 0
		}(),
	}, nil
}
