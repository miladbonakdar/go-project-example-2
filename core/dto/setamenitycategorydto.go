package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SetAmenityCategoryDto struct {
	AmenityId         int  `json:"amenityId"`
	AmenityCategoryId uint `json:"amenityCategoryId"`
}

func (a SetAmenityCategoryDto) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.AmenityId, validation.Required, validation.Min(0)),
		validation.Field(&a.AmenityCategoryId, validation.Required),
	)
}
