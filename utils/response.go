package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

func JSON(write http.ResponseWriter, status int, payload Response) {
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(status)
	_ = json.NewEncoder(write).Encode(payload)
}

func Success(write http.ResponseWriter, data any, message string) {
	JSON(write, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(write http.ResponseWriter, status int, message string, errs any) {
	JSON(write, status, Response{
		Success: false,
		Message: message,
		Errors:  errs,
	})
}

type Pagination struct {
	Total       int  `json:"total"`
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	TotalPages  int  `json:"total_pages"`
	HasPrevPage bool `json:"has_prev_page"`
	HasNextPage bool `json:"has_next_page"`
}

func SuccessWithPagination(write http.ResponseWriter, data any, message string, meta Pagination) {
	JSON(write, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}
