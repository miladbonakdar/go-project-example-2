package dto

import "hotel-engine/utils/indraframework"

type OrderRefundResponseDto struct {
	OrderId         string                         `json:"orderId"`
	RefundRequestId int64                          `json:"refundRequestId"`
	Error           *indraframework.IndraException `json:"error"`
}

func (a *OrderRefundResponseDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
