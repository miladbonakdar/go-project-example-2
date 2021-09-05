package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SyncSomeHotelsDto struct {
	HotelIds []string                       `json:"hotel_ids"`
}

func (a SyncSomeHotelsDto) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.HotelIds,
			validation.Each(validation.Length(1, 256))),
	)
}
