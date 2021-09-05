package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SetAmenityIconDto struct {
	AmenityId int    `json:"amenityId"`
	IconUrl   string `json:"iconUrl"`
}

func (a SetAmenityIconDto) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.AmenityId, validation.Required, validation.Min(0)),
		validation.Field(&a.IconUrl, validation.Required),
	)
}
