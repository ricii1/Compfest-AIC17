package controller

import (
	"net/http"

	"github.com/Caknoooo/go-gin-clean-starter/dto"
	"github.com/Caknoooo/go-gin-clean-starter/service"
	"github.com/Caknoooo/go-gin-clean-starter/utils"
	"github.com/gin-gonic/gin"
)

type (
	ReportController interface {
		CreateReport(ctx *gin.Context)
	}

	reportController struct {
		reportService service.ReportService
	}
)

func NewReportController(rs service.ReportService) ReportController {
	return &reportController{
		reportService: rs,
	}
}

func (c *reportController) CreateReport(ctx *gin.Context) {
	var report dto.CreateReportRequest
	if err := ctx.ShouldBind(&report); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.reportService.CreateReport(ctx, report)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_SEND_REPORT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_SEND_REPORT, result)
	ctx.JSON(http.StatusOK, res)
}
