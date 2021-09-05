package dto

import "hotel-engine/utils/indraframework"

type OrderPayByAccountResponseDto struct {
	TransactionStatus string                         `json:"transactionStatus"`
	RequestId         string                         `json:"requestId"`
	TransactionIds    []string                       `json:"transactionIds"`
	ResultMessage     string                         `json:"resultMessage"`
	Error             *indraframework.IndraException `json:"error"`
}

func (a *OrderPayByAccountResponseDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
