package repository

import (
	"context"
	"fmt"

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
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&report).Error; err != nil {
		fmt.Println(err)
		return entity.Report{}, err
	}

	return report, nil
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
