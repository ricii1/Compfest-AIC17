package migrations

import (
	"github.com/Caknoooo/go-gin-clean-starter/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.RefreshToken{},
		&entity.Report{},
		&entity.Tag{},
	); err != nil {
		return err
	}

	return nil
}
