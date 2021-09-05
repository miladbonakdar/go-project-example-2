package dto

import "hotel-engine/utils/indraframework"

type HotelsPageResponseDto struct {
	PageNumber int                            `json:"pageNumber"`
	PageSize   int                            `json:"pageSize"`
	Total      int                            `json:"total"`
	Hotels     []HotelDto                     `json:"hotels"`
	Error      *indraframework.IndraException `json:"error"`
}

func (a *HotelsPageResponseDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
