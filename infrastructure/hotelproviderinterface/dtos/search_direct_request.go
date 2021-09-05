package dtos

import (
	"hotel-engine/utils/date"
	"time"
)

func NewSearchDirectRequest(hotelId string, checkIn, checkOut time.Time) ProviderRequestSessionDto {
	return ProviderRequestSessionDto{
		CheckIn:  checkIn.Format(date.LayoutISO),
		CheckOut: checkOut.Format(date.LayoutISO),
		Rooms: []ProviderRequestSessionRoomDto{
			{
				Adults:   []int{30},
				Children: []int{},
			},
		},
		Destination: ProviderRequestSessionDestinationDto{
			Id:   hotelId,
			Type: "Hotel",
		},
	}
}
