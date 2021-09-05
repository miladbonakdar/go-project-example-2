package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"hotel-engine/core/common"
	"hotel-engine/utils/date"
)

type HotelDetailsDto struct {
	HotelId          string           `json:"hotelId"`
	CheckIn          string           `json:"checkIn"`
	CheckOut         string           `json:"checkOut"`
	Rooms            []RequestRoomDto `json:"rooms"`
	SkipLoadingRooms bool             `json:"skipLoadingRooms,omitempty"`
}

type RequestRoomDto struct {
	Adults   []int `json:"adults"`
	Children []int `json:"children"`
}

func (a RequestRoomDto) Validate() error {
	if len(a.Adults) == 0 {
		a.Adults = []int{30}
	}
	return validation.ValidateStruct(&a,
		validation.Field(&a.Adults, validation.Each(validation.Min(common.MaxAgeAsAChild))),
		validation.Field(&a.Children, validation.Each(validation.Max(common.MaxAgeAsAChild-1),
			validation.Min(0))),
	)
}

func (a HotelDetailsDto) Validate() error {
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

func (a HotelDetailsDto) CalculateHours() int {
	checkIn := date.StringToDateOrDefault(a.CheckIn)
	checkOut := date.StringToDateOrDefault(a.CheckOut)
	return int(checkOut.Sub(checkIn).Hours())
}

func (a HotelDetailsDto) CanSkipOptions() bool {
	standardAdults := len(a.Rooms) == 1 && len(a.Rooms[0].Adults) == 1
	return a.SkipLoadingRooms && a.CalculateHours() <= 24 && standardAdults
}

func (a *HotelDetailsDto) SetDefaults() {
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
