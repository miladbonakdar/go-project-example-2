package dto

import "hotel-engine/utils/indraframework"

type TokenDto struct {
	Token string                         `json:"token"`
	Error *indraframework.IndraException `json:"error"`
}

func NewTokenDto(token string) *TokenDto {
	return &TokenDto{
		Token: token,
	}
}

func (a *TokenDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
