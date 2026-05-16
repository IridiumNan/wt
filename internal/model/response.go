package model

import (
	"net/http"
)

type APIResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func SuccessfulResponse(data any, message string) *APIResponse {
	if message == "" {
		message = "success"
	}

	return &APIResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func BadRequestResponse(errorMsg string) *APIResponse {
	return &APIResponse{
		Code:   http.StatusBadRequest,
		Status: "bad_request",
		Error:  errorMsg,
	}
}

func UnauthorizedResponse(errorMsg string) *APIResponse {
	return &APIResponse{
		Code:   http.StatusUnauthorized,
		Status: "unauthorized",
		Error:  errorMsg,
	}
}

func ForbiddenResponse(errorMsg string) *APIResponse {
	return &APIResponse{
		Code:   http.StatusForbidden,
		Status: "forbidden",
		Error:  errorMsg,
	}
}

func NotFoundResponse(errorMsg string) *APIResponse {
	return &APIResponse{
		Code:   http.StatusNotFound,
		Status: "not_found",
		Error:  errorMsg,
	}
}

func InternalErrorResponse(errorMsg string) *APIResponse {
	return &APIResponse{
		Code:   http.StatusInternalServerError,
		Status: "internal_error",
		Error:  errorMsg,
	}
}
