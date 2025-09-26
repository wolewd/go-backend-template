package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type Pagination struct {
	Total       int  `json:"total"`
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	TotalPages  int  `json:"total_pages"`
	HasPrevPage bool `json:"has_prev_page"`
	HasNextPage bool `json:"has_next_page"`
}

func JSON(w http.ResponseWriter, status int, payload Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func Success(w http.ResponseWriter, data interface{}, message string) {
	JSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithPagination(w http.ResponseWriter, data interface{}, message string, meta Pagination) {
	JSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(w http.ResponseWriter, status int, message string, errs interface{}) {
	JSON(w, status, Response{
		Success: false,
		Message: message,
		Errors:  errs,
	})
}
