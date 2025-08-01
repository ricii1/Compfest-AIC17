package dto

import (
	"errors"
	"mime/multipart"
)

const (
	// Failed
	MESSAGE_FAILED_SEND_REPORT = "gagal mengirim laporan"

	// Success
	MESSAGE_SUCCESS_SEND_REPORT = "berhasil mengirim laporan"
)

var (
	ErrInvalidImageType      = errors.New("invalid image type")
	ErrImageSizeTooLarge     = errors.New("image size too large")
	ErrEmptyFileUploaded     = errors.New("empty file uploaded")
	ErrInvalidImageExtension = errors.New("invalid image extension")
	ErrFailedToOpenImageFile = errors.New("failed to open image file")
	ErrCreateReport          = errors.New("failed to create report")
	ErrEmptyContent          = errors.New("error empty content")

// ErrCreateUser             = errors.New("failed to create user")
)

type (
	CreateReportRequest struct {
		Text  string                `json:"text" form:"text"`
		Image *multipart.FileHeader `json:"image" form:"image"`
	}
	CreateReportResponse struct {
		ID    string `json:"id"`
		Text  string `json:"text"`
		Image string `json:"image"`
	}
)
