package dto

import "hotel-engine/utils/indraframework"

type FinalizeOrderResponseDto struct {
	PaymentResult *OrderPayByAccountResponseDto  `json:"paymentResult"`
	StatusResult  *OrderStatusResponseDto        `json:"statusResult"`
	Error         *indraframework.IndraException `json:"error"`
}

func (a *FinalizeOrderResponseDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
