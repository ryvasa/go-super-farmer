package helper

import (
	"encoding/json"
	"net/http"
)

// ResponseBody adalah struktur umum untuk respons API
type ResponseBodyPagination[T any] struct {
	Status  int           `json:"status"`
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    Pagination[T] `json:"data"`
	Errors  interface{}   `json:"errors,omitempty"`
}

// Pagination adalah struktur umum untuk paginasi
type Pagination[T any] struct {
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Data       []T   `json:"data"`
}

type ResponseBody[T any] struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    T           `json:"data"`
	Errors  interface{} `json:"errors"`
}

func ParseResponseBody[T any](resp *http.Response) (*ResponseBody[T], error) {
	var responseBody ResponseBody[T]
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&responseBody)
	if err != nil {
		return nil, err
	}
	return &responseBody, nil
}

func ParseResponseBodyPagination[T any](resp *http.Response) (*ResponseBodyPagination[T], error) {
	var responseBody ResponseBodyPagination[T]
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&responseBody)
	if err != nil {
		return nil, err
	}
	return &responseBody, nil
}
