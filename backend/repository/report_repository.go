package repository

import (
	"context"
	"fmt"

	"github.com/Caknoooo/go-gin-clean-starter/entity"
	"gorm.io/gorm"
)

type (
	ReportRepository interface {
		CreateReport(ctx context.Context, tx *gorm.DB, report entity.Report) (entity.Report, error)
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
