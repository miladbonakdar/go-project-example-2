package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SetHotelSeoRequestDto struct {
	HotelId    string      `json:"hotelId"`
	SeoDetails HotelSeoDto `json:"seoDetails"`
}

func (a SetHotelSeoRequestDto) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.HotelId, validation.Required),
	)
}
