package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type OrderRefundRequestDto struct {
	OrderId         string `json:"orderId"`
	JabamaOrderID   int64  `json:"jabamaOrderId"`
	RefundRequestID int64  `json:"refundRequestId"`
}

func (a OrderRefundRequestDto) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.OrderId, validation.Required),
		validation.Field(&a.JabamaOrderID, validation.Required),
		validation.Field(&a.RefundRequestID, validation.Required),
	)
}
