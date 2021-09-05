package dtos

import "hotel-engine/core/dto"

type DirectRequestDto struct {
	Destination DirectRequestDestinationDto `json:"destination"`
	CheckIn     string                      `json:"checkIn"`
	CheckOut    string                      `json:"checkOut"`
	Rooms       []DirectRequestRoomDto      `json:"rooms"`
}

type DirectRequestDestinationDto struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type DirectRequestRoomDto struct {
	Adults   []int `json:"adults"`
	Children []int `json:"children"`
}

func NewDirectRequestDto(hotelId string, checkIn, checkOut string, rooms []dto.RequestRoomDto) *DirectRequestDto {
	directRooms := make([]DirectRequestRoomDto, 0, len(rooms))
	for _, room := range rooms {
		directRooms = append(directRooms, DirectRequestRoomDto{
			Adults:   room.Adults,
			Children: room.Children,
		})
	}

	return &DirectRequestDto{
		Destination: DirectRequestDestinationDto{
			Id:   hotelId,
			Type: "Hotel",
		},
		CheckIn:  checkIn,
		CheckOut: checkOut,
		Rooms:    directRooms,
	}
}
