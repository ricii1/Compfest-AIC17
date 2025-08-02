package provider

import (
	"github.com/Caknoooo/go-gin-clean-starter/controller"
	"github.com/Caknoooo/go-gin-clean-starter/repository"
	"github.com/Caknoooo/go-gin-clean-starter/service"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func ProvideReportDependencies(injector *do.Injector, db *gorm.DB, jwtService service.JWTService) {
	// Repository
	reportRepository := repository.NewReportRepository(db)
	userRepository := repository.NewUserRepository(db)
	// Service
	reportService := service.NewReportService(userRepository, reportRepository, db)
	// userService := service.NewUserService(userRepository, refreshTokenRepository, jwtService, db)

	// Controller
	do.Provide(
		injector, func(i *do.Injector) (controller.ReportController, error) {
			return controller.NewReportController(reportService), nil
		},
	)
}
