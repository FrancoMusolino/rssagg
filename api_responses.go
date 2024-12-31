package main

type ApiResponse interface {
	IsSuccess() bool
	GetMessage() string
}

type ApiSuccessResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func NewApiSuccessResponse[T any](message string, data T) ApiSuccessResponse[T] {
	return ApiSuccessResponse[T]{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func (r ApiSuccessResponse[T]) IsSuccess() bool {
	return r.Success
}

func (r ApiSuccessResponse[T]) GetMessage() string {
	return r.Message
}

type ApiFailureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewApiFailureResponse(message string) ApiFailureResponse {
	return ApiFailureResponse{
		Success: false,
		Message: message,
	}
}

func (r ApiFailureResponse) IsSuccess() bool {
	return r.Success
}

func (r ApiFailureResponse) GetMessage() string {
	return r.Message
}
