package dto

import (
	"hotel-engine/utils/indraframework"
)

type HotelPDPDto struct {
	HotelDto

	Rooms []RequestRoomDto `json:"requestedRooms"`
}

func (a *HotelPDPDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}

func (a *HotelPDPDto) SetProvince(province string) {
	a.Province = province
}

func (a *HotelPDPDto) SetType(hotelType string) {
	a.Type = hotelType
}
