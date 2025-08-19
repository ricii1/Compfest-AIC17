package routes

import (
	"github.com/Caknoooo/go-gin-clean-starter/constants"
	"github.com/Caknoooo/go-gin-clean-starter/controller"
	"github.com/Caknoooo/go-gin-clean-starter/middleware"
	"github.com/Caknoooo/go-gin-clean-starter/service"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func Reports(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	reportController := do.MustInvoke[controller.ReportController](injector)

	routes := route.Group("/api/reports")
	{
		// Reports
		routes.POST("", middleware.Authenticate(jwtService), reportController.CreateReport)
		routes.GET("", middleware.Authenticate(jwtService), reportController.GetAllReports)
		routes.GET("/:id", middleware.Authenticate(jwtService), reportController.GetReportById)
		routes.GET("/user/:id", middleware.Authenticate(jwtService), reportController.GetReportsByUserId)
		routes.POST("/:id/status", middleware.Authenticate(jwtService), reportController.UpdateReportStatus)
	}
}
