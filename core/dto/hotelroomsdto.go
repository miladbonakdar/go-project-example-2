package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type HotelRoomsDto struct {
	HotelId  string           `json:"hotelId"`
	CheckIn  string           `json:"checkIn"`
	CheckOut string           `json:"checkOut"`
	Rooms    []RequestRoomDto `json:"rooms"`
}
type HotelRoomsWithSessionDto struct {
	HotelId   string `json:"hotelId"`
	SessionId string `json:"sessionId"`
}

func (a HotelRoomsDto) Validate() error {
	if err := CheckForDate(a.CheckIn); err != nil {
		return err
	}
	if err := CheckForDate(a.CheckOut); err != nil {
		return err
	}
	for _, room := range a.Rooms {
		if err := room.Validate(); err != nil {
			return err
		}
	}
	return validation.ValidateStruct(&a,
		validation.Field(&a.CheckOut, validation.Required),
		validation.Field(&a.CheckIn, validation.Required),
		validation.Field(&a.HotelId, validation.Required),
	)
}

func (a *HotelRoomsDto) SetDefaults() {
	if a.Rooms == nil || len(a.Rooms) == 0 {
		a.Rooms = []RequestRoomDto{
			{
				Adults:   []int{30},
				Children: []int{},
			},
		}
	}
	if a.Rooms[0].Adults == nil || len(a.Rooms[0].Adults) == 0 {
		a.Rooms[0].Adults = []int{30}
	}
}
