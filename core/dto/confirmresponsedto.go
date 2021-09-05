package dto

import "hotel-engine/utils/indraframework"

type ConfirmResponseDto struct {
	OrderId string                         `json:"orderId"`
	Error   *indraframework.IndraException `json:"error"`
}

func (a *ConfirmResponseDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
