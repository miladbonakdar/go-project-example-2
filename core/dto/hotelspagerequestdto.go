package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type HotelsPageRequestDto struct {
	PageNumber int    `json:"pageNumber"`
	PageSize   int    `json:"pageSize"`
	Search     string `json:"search"`
}

func (a HotelsPageRequestDto) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.PageNumber, validation.Required, validation.Min(1)),
		validation.Field(&a.PageSize, validation.Required, validation.Min(1)),
	)
}
