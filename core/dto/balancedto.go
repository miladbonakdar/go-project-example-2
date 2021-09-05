package dto

import "hotel-engine/utils/indraframework"

type BalanceDto struct {
	Balance float64                        `json:"balance"`
	Error   *indraframework.IndraException `json:"error"`
}

func (a *BalanceDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}

func NewBalanceDto(balance float64) BalanceDto {
	return BalanceDto{
		Balance: balance,
		Error:   nil,
	}
}
