package service

import (
	"context"
	"fmt"
	"os"

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
		GetReportsByUserId(ctx context.Context, userId string, req dto.PaginationRequest) (dto.ReportPaginationResponse, error)
		UpdateReportStatus(ctx context.Context, reportId string, status entity.ReportStatus) (dto.UpdateStatusReportResponse, error)
		CountReportStatus(ctx context.Context) (dto.CountReportResponse, error)
		GetReportsByStatus(ctx context.Context, status string, req dto.PaginationRequest) (dto.ReportPaginationResponse, error)
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

func (s *reportService) GetReportsByUserId(ctx context.Context, userId string, req dto.PaginationRequest) (dto.ReportPaginationResponse, error) {
	reports, err := s.reportRepo.GetReportsByUserId(ctx, nil, userId, req)
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

func (s *reportService) UpdateReportStatus(ctx context.Context, reportId string, status entity.ReportStatus) (dto.UpdateStatusReportResponse, error) {
	if status != entity.StatusUnverified && status != entity.StatusVerified && status != entity.StatusRejected && status != entity.StatusCompleted && status != entity.StatusHandled {
		return dto.UpdateStatusReportResponse{}, dto.ErrUpdateReportStatus
	}
	reports, err := s.reportRepo.UpdateReportStatus(ctx, nil, reportId, status)
	if err != nil {
		return dto.UpdateStatusReportResponse{}, dto.ErrUpdateReportStatus
	}

	return reports, nil
}

func (s *reportService) CountReportStatus(ctx context.Context) (dto.CountReportResponse, error) {
	counts, err := s.reportRepo.CountReportStatus(ctx, nil)
	if err != nil {
		return dto.CountReportResponse{}, dto.ErrGetReports
	}
	return dto.CountReportResponse{
		Total:      counts.Total,
		Unverified: counts.Unverified,
		Verified:   counts.Verified,
		Rejected:   counts.Rejected,
		Handled:    counts.Handled,
		Completed:  counts.Completed,
	}, nil
}

func (s *reportService) GetReportsByStatus(ctx context.Context, status string, req dto.PaginationRequest) (dto.ReportPaginationResponse, error) {
	reportStatus := entity.ReportStatus(status)
	if reportStatus != entity.StatusUnverified && reportStatus != entity.StatusVerified && reportStatus != entity.StatusRejected && reportStatus != entity.StatusCompleted && reportStatus != entity.StatusHandled {
		return dto.ReportPaginationResponse{}, dto.ErrGetReports
	}
	reports, err := s.reportRepo.GetReportsByStatus(ctx, nil, reportStatus, req)
	if err != nil {
		return dto.ReportPaginationResponse{}, dto.ErrGetReports
	}

	var datas []dto.ReportResponse
	for _, report := range reports.Reports {
		user, err := s.userRepo.GetUserById(ctx, nil, report.UserID)
		if err != nil {
			return dto.ReportPaginationResponse{}, dto.ErrGetUserById
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
