package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SetHotelFaqRequestDto struct {
	HotelId  string        `json:"hotelId"`
	FAQTitle string        `json:"faqTitle"`
	FAQList  []HotelFAQDto `json:"faqList"`
}

type HotelFAQDto struct {
	ID       uint   `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func (a SetHotelFaqRequestDto) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.HotelId, validation.Required),
	)
}
