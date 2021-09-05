package dto

import (
	"hotel-engine/utils/indraframework"
)

//FindHotelResponseDto for extracting hetol data
type FindHotelResponseDto struct {
	Hotels *[]HotelDto                    `json:"hotels"`
	Error  *indraframework.IndraException `json:"error"`
}

func (a *FindHotelResponseDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
