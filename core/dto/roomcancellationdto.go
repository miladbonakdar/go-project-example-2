package dto

import "hotel-engine/utils/indraframework"

type RoomCancellationPolicyDto struct {
	NonRefundable   bool                           `json:"nonRefundable"`
	GeneralPolicies []string                       `json:"GeneralPolicies"`
	Error           *indraframework.IndraException `json:"error"`
}

func (a *RoomCancellationPolicyDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
