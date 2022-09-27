package web

import "github.com/go-playground/validator"

const (
	SUCCESS = "success"
	FAIL    = "fail"
)

type (
	CustomsAppendToICP struct {
		FileName   string   `json:"file_name" validate:"required"`
		CustomsIds []string `json:"customs_ids" validate:"required"`
	}

	CustomValidator struct {
		Validator *validator.Validate
	}

	IcpResponse struct {
		Status   string   `json:"status"`
		FileName string   `json:"file_name"`
		Errors   []string `json:"errors"`
	}
)
