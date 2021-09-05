package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type OptionInfoRequestDto struct {
	OptionId  string `json:"optionId"`
	HotelId   string `json:"hotelId"`
	SessionId string `json:"sessionId"`
}

func (a OptionInfoRequestDto) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.OptionId, validation.Required),
		validation.Field(&a.SessionId, validation.Required),
		validation.Field(&a.HotelId, validation.Required),
	)
}
