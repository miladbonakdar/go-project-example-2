package dto

import "hotel-engine/utils/indraframework"

type AvailableResponseDto struct {
	OrderId      string                         `json:"orderId"`
	OptionId     string                         `json:"optionId"`
	TotalPrice   int64                          `json:"totalPrice"`
	IndraOrderId int64                          `json:"IndraOrderId"`
	Status       string                         `json:"status"`
	CheckIn      string                         `json:"checkIn"`
	CheckOut     string                         `json:"checkOut"`
	Error        *indraframework.IndraException `json:"error"`
}

func (a *AvailableResponseDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
