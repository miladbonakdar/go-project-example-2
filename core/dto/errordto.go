package dto

import "net/http"

//error dto with helpers function for
type ErrorDto struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewErrorDto(message string, code int) *ErrorDto {
	return &ErrorDto{
		Message: message,
		Code:    code,
	}
}

func NewNotFound(message string) *ErrorDto {
	return NewErrorDto(message, http.StatusNotFound)
}

func NewInternalServerError(message string) *ErrorDto {
	return NewErrorDto(message, http.StatusInternalServerError)
}

func NewBadRequest(message string) *ErrorDto {
	return NewErrorDto(message, http.StatusBadRequest)
}

func NotFoundFromError(err error) *ErrorDto {
	return NewErrorDto(err.Error(), http.StatusNotFound)
}

func InternalServerErrorFrom(err error) *ErrorDto {
	return NewErrorDto(err.Error(), http.StatusInternalServerError)
}

func BadRequestFromError(err error) *ErrorDto {
	return NewErrorDto(err.Error(), http.StatusBadRequest)
}
