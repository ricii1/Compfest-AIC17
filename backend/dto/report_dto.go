package dto

import (
	"errors"
	"mime/multipart"

	"github.com/Caknoooo/go-gin-clean-starter/entity"
)

const (
	// Failed
	MESSAGE_FAILED_SEND_REPORT            = "gagal mengirim laporan"
	MESSAGE_FAILED_GET_REPORTS            = "gagal mendapatkan laporan"
	MESSAGE_FAILED_GET_REPORT_BY_ID       = "gagal mendapatkan laporan dari id"
	MESSAGE_FAILED_GET_REPORTS_BY_USER_ID = "gagal mendapatkan laporan berdasarkan user id"
	MESSAGE_FAILED_DENIED                 = "akses ditolak"

	// Success
	MESSAGE_SUCCESS_SEND_REPORT            = "berhasil mengirim laporan"
	MESSAGE_SUCCESS_GET_REPORTS            = "berhasil mendapatkan laporan"
	MESSAGE_SUCCESS_GET_REPORT_BY_ID       = "berhasil mendapatkan laporan dari id"
	MESSAGE_SUCCESS_GET_REPORTS_BY_USER_ID = "berhasil mendapatkan laporan berdasarkan user id"
)

var (
	ErrInvalidImageType      = errors.New("tipe gambar tidak valid")
	ErrImageSizeTooLarge     = errors.New("ukuran gambar terlalu besar")
	ErrEmptyFileUploaded     = errors.New("file kosong")
	ErrInvalidImageExtension = errors.New("extensi gambar tidak valid")
	ErrFailedToOpenImageFile = errors.New("gagal membuka gambar")
	ErrCreateReport          = errors.New("gagal membuat laporan")
	ErrEmptyContent          = errors.New("konten kosong")
	ErrGetReports            = errors.New("gagal mendapatkan laporan")
	ErrGetReportById         = errors.New("gagal mendapatkan laporan dari id")
	ErrUpdateReportStatus    = errors.New("gagal memperbarui status laporan")
	ErrUpdateReportInference = errors.New("gagal memperbarui inferensi laporan")

// ErrCreateUser             = errors.New("failed to create user")
)

type (
	CreateReportRequest struct {
		Text     string                `json:"text" form:"text"`
		Location string                `json:"location" form:"location"`
		Image    *multipart.FileHeader `json:"image" form:"image"`
	}
	CreateReportResponse struct {
		ID       string `json:"id"`
		Text     string `json:"text"`
		Image    string `json:"image"`
		Location string `json:"location"`
	}

	ReportResponse struct {
		ID             string `json:"id"`
		Text           string `json:"text"`
		Image          string `json:"image"`
		Location       string `json:"location"`
		Status         string `json:"status"`
		Upvotes        int    `json:"upvotes"`
		ShareCount     int    `json:"share_count"`
		TagID          string `json:"tag_id"`
		UserID         string `json:"user_id"`
		Username       string `json:"username"`
		PredConfidence int    `json:"pred_confidence"`
		CreatedAt      string `json:"created_at"`
	}

	ReportPaginationResponse struct {
		Data []ReportResponse `json:"data"`
		PaginationResponse
	}

	GetAllReportResponse struct {
		Reports []entity.Report `json:"reports"`
		PaginationResponse
	}

	UpdateStatusReportRequest struct {
		Status entity.ReportStatus `json:"status"`
	}

	UpdateStatusReportResponse struct {
		ID     string              `json:"id"`
		Status entity.ReportStatus `json:"status"`
	}

	StatusCount struct {
		Status string
		Count  int64
	}
	CountReportResponse struct {
		Total      int64 `json:"total"`
		Unverified int64 `json:"unverified"`
		Verified   int64 `json:"verified"`
		Rejected   int64 `json:"rejected"`
		Handled    int64 `json:"handled"`
		Completed  int64 `json:"completed"`
	}
	InferenceRequest struct {
		ReportID string `json:"report_id"`
		Class    string `json:"class"`
		Location string `json:"location"`
	}
	InferenceTag struct {
		TagID    string `json:"tag_id"`
		Class    string `json:"class"`
		Location string `json:"location"`
	}
	InferenceResponse struct {
		Data []InferenceTag `json:"data"`
	}
)
