package repository

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub/v2"
	"github.com/Caknoooo/go-gin-clean-starter/dto"
	"github.com/Caknoooo/go-gin-clean-starter/entity"
	"gorm.io/gorm"
)

type (
	ReportRepository interface {
		CreateReport(ctx context.Context, tx *gorm.DB, report entity.Report) (entity.Report, error)
		GetAllReportsWithPagination(
			ctx context.Context,
			tx *gorm.DB,
			req dto.PaginationRequest,
		) (dto.GetAllReportResponse, error)
		GetReportById(ctx context.Context, tx *gorm.DB, reportId string) (entity.Report, error)
		GetReportsByUserId(ctx context.Context, tx *gorm.DB, userId string, req dto.PaginationRequest) (dto.GetAllReportResponse, error)
		UpdateReportStatus(ctx context.Context, tx *gorm.DB, reportId string, status entity.ReportStatus) (dto.UpdateStatusReportResponse, error)
		CountReportStatus(ctx context.Context, tx *gorm.DB) (dto.CountReportResponse, error)
		GetReportsByStatus(ctx context.Context, tx *gorm.DB, status entity.ReportStatus, req dto.PaginationRequest) (dto.GetAllReportResponse, error)
		UpdateReportInference(ctx context.Context, tx *gorm.DB, report entity.Report, class string, location string) (entity.Tag, error)
	}

	reportRepository struct {
		db *gorm.DB
	}
)

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{
		db: db,
	}
}

func (r *reportRepository) CreateReport(ctx context.Context, tx *gorm.DB, report entity.Report) (entity.Report, error) {
	var db *gorm.DB
	if tx == nil {
		db = r.db
	} else {
		db = tx
	}

	var createdReport entity.Report
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&report).Error; err != nil {
			return err
		}
		gcp_project_id := os.Getenv("GCP_PROJECT_ID")
		topic_id := os.Getenv("GCP_TOPIC_ID")
		client, err := pubsub.NewClient(ctx, gcp_project_id)
		if err != nil {
			return err
		}
		defer client.Close()

		publisher := client.Publisher(topic_id)
		result := publisher.Publish(ctx, &pubsub.Message{
			Data: []byte(fmt.Sprintf(`{"report_id": "%s", "text": "%s", "image_url": "%s"}`, report.ID, report.Text, report.Image)),
		})
		_, err = result.Get(ctx)
		if err != nil {
			return err
		}
		createdReport = report
		return nil
	})

	if err != nil {
		return entity.Report{}, err
	}
	return createdReport, nil
}

func (r *reportRepository) GetAllReportsWithPagination(
	ctx context.Context,
	tx *gorm.DB,
	req dto.PaginationRequest,
) (dto.GetAllReportResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var reports []entity.Report
	var err error
	var count int64

	req.Default()

	query := tx.WithContext(ctx).Model(&entity.Report{})
	if req.Search != "" {
		query = query.Where("text LIKE ?", "%"+req.Search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.GetAllReportResponse{}, err
	}

	if err := query.Scopes(Paginate(req)).Find(&reports).Error; err != nil {
		return dto.GetAllReportResponse{}, err
	}

	totalPage := TotalPage(count, int64(req.PerPage))
	return dto.GetAllReportResponse{
		Reports: reports,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, err
}

func (r *reportRepository) GetReportById(ctx context.Context, tx *gorm.DB, reportId string) (entity.Report, error) {
	if tx == nil {
		tx = r.db
	}

	var report entity.Report
	if err := tx.WithContext(ctx).First(&report, "id = ?", reportId).Error; err != nil {
		return entity.Report{}, err
	}

	return report, nil
}

func (r *reportRepository) GetReportsByUserId(ctx context.Context, tx *gorm.DB, userId string, req dto.PaginationRequest) (dto.GetAllReportResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var reports []entity.Report
	var err error
	var count int64

	req.Default()

	query := tx.WithContext(ctx).Model(&entity.Report{})
	query = query.Where("user_id = ?", userId)
	if req.Search != "" {
		query = query.Where("text LIKE ?", "%"+req.Search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.GetAllReportResponse{}, err
	}

	if err := query.Scopes(Paginate(req)).Find(&reports).Error; err != nil {
		return dto.GetAllReportResponse{}, err
	}

	totalPage := TotalPage(count, int64(req.PerPage))
	return dto.GetAllReportResponse{
		Reports: reports,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, err
}

func (r *reportRepository) UpdateReportStatus(ctx context.Context, tx *gorm.DB, reportId string, status entity.ReportStatus) (dto.UpdateStatusReportResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var report entity.Report
	if err := tx.WithContext(ctx).First(&report, "id = ?", reportId).Error; err != nil {
		return dto.UpdateStatusReportResponse{}, err
	}

	report.Status = status
	if err := tx.WithContext(ctx).Save(&report).Error; err != nil {
		return dto.UpdateStatusReportResponse{}, err
	}

	return dto.UpdateStatusReportResponse{
		ID:     report.ID,
		Status: report.Status,
	}, nil
}

func (r *reportRepository) CountReportStatus(ctx context.Context, tx *gorm.DB) (dto.CountReportResponse, error) {
	if tx == nil {
		tx = r.db
	}
	var counts []dto.StatusCount
	if err := tx.WithContext(ctx).Model(&entity.Report{}).Select("status, COUNT(*) as count").Group("status").Scan(&counts).Error; err != nil {
		return dto.CountReportResponse{}, err
	}
	var amount dto.CountReportResponse
	for _, count := range counts {
		switch count.Status {
		case "unverified":
			amount.Unverified = count.Count
		case "verified":
			amount.Verified = count.Count
		case "rejected":
			amount.Rejected = count.Count
		case "handled":
			amount.Handled = count.Count
		case "completed":
			amount.Completed = count.Count
		}
		amount.Total += count.Count
	}
	return amount, nil
}

func (r *reportRepository) GetReportsByStatus(ctx context.Context, tx *gorm.DB, status entity.ReportStatus, req dto.PaginationRequest) (dto.GetAllReportResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var reports []entity.Report
	var err error
	var count int64

	req.Default()

	query := tx.WithContext(ctx).Model(&entity.Report{})
	query = query.Where("status = ?", status)
	if req.Search != "" {
		query = query.Where("text LIKE ?", "%"+req.Search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.GetAllReportResponse{}, err
	}

	if err := query.Scopes(Paginate(req)).Find(&reports).Error; err != nil {
		return dto.GetAllReportResponse{}, err
	}

	totalPage := TotalPage(count, int64(req.PerPage))
	return dto.GetAllReportResponse{
		Reports: reports,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, err
}

func (r *reportRepository) UpdateReportInference(ctx context.Context, tx *gorm.DB, report entity.Report, class string, location string) (entity.Tag, error) {
	if tx == nil {
		tx = r.db
	}

	var tag entity.Tag
	tag.Reports = append(tag.Reports, report)
	tag.Class = class
	if location == "" {
		tag.Location = report.Location
	} else {
		tag.Location = location
	}

	if err := tx.WithContext(ctx).Save(&tag).Error; err != nil {
		return entity.Tag{}, err
	}

	return tag, nil
}
