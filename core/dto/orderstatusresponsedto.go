package dto

import "hotel-engine/utils/indraframework"

type OrderStatusResponseDto struct {
	OrderId string                         `json:"orderId"`
	Status  string                         `json:"status"`
	Error   *indraframework.IndraException `json:"error"`
}

func (a *OrderStatusResponseDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
