package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type FinalizeOrderDto struct {
	OrderId string `json:"orderId"`
}

func (a FinalizeOrderDto) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.OrderId, validation.Required),
	)
}
