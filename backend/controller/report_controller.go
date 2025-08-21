package controller

import (
	"fmt"
	"net/http"

	"github.com/Caknoooo/go-gin-clean-starter/dto"
	"github.com/Caknoooo/go-gin-clean-starter/service"
	"github.com/Caknoooo/go-gin-clean-starter/utils"
	"github.com/gin-gonic/gin"
)

type (
	ReportController interface {
		CreateReport(ctx *gin.Context)
		GetAllReports(ctx *gin.Context)
		GetReportById(ctx *gin.Context)
		GetReportsByUserId(ctx *gin.Context)
		UpdateReportStatus(ctx *gin.Context)
		CountReportStatus(ctx *gin.Context)
		GetReportsByStatus(ctx *gin.Context)
		InferenceStatus(ctx *gin.Context)
	}

	reportController struct {
		reportService service.ReportService
		userService   service.UserService
	}
)

func NewReportController(rs service.ReportService, us service.UserService) ReportController {
	return &reportController{
		reportService: rs,
		userService:   us,
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

func (c *reportController) GetAllReports(ctx *gin.Context) {
	var req dto.PaginationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	reports, err := c.reportService.GetAllReports(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_REPORTS, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_REPORTS, reports)
	ctx.JSON(http.StatusOK, res)
}

func (c *reportController) GetReportById(ctx *gin.Context) {
	reportId := ctx.Param("id")
	result, err := c.reportService.GetReportById(ctx.Request.Context(), reportId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_REPORT_BY_ID, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_REPORT_BY_ID, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *reportController) GetReportsByUserId(ctx *gin.Context) {
	userId := ctx.Param("id")
	var req dto.PaginationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	result, err := c.reportService.GetReportsByUserId(ctx.Request.Context(), userId, req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_REPORTS_BY_USER_ID, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_REPORTS_BY_USER_ID, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *reportController) UpdateReportStatus(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)
	user, err := c.userService.GetUserById(ctx.Request.Context(), userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	if user.Role != "admin" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DENIED_ACCESS, dto.MESSAGE_FAILED_DENIED, nil)
		ctx.JSON(http.StatusForbidden, res)
		return
	}
	reportId := ctx.Param("id")
	var req dto.UpdateStatusReportRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	result, err := c.reportService.UpdateReportStatus(ctx.Request.Context(), reportId, req.Status)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_REPORT_BY_ID, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_REPORT_BY_ID, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *reportController) CountReportStatus(ctx *gin.Context) {
	result, err := c.reportService.CountReportStatus(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_REPORTS, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_REPORTS, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *reportController) GetReportsByStatus(ctx *gin.Context) {
	var req dto.PaginationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	status := ctx.Param("status")
	result, err := c.reportService.GetReportsByStatus(ctx.Request.Context(), status, req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_REPORTS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_REPORTS, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *reportController) InferenceStatus(ctx *gin.Context) {
	fmt.Println("TES")
	var req dto.InferenceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	token := ctx.GetHeader("X-WEBHOOK-TOKEN")
	result, err := c.reportService.InferenceStatus(ctx.Request.Context(), req, token)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_REPORTS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_REPORTS, result)
	ctx.JSON(http.StatusOK, res)
}
