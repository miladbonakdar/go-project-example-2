package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"hotel-engine/core/common"
)

type AvailableDto struct {
	SessionId   string             `json:"sessionId"`
	HotelId     string             `json:"hotelId"`
	OptionId    string             `json:"optionId"`
	NationalId  string             `json:"nationalId"`
	PhoneNumber string             `json:"phoneNumber"`
	LateCheckIn string             `json:"lateCheckIn,omitempty"`
	CheckIn     string             `json:"checkIn"`
	CheckOut    string             `json:"checkOut"`
	Rooms       []AvailableRoomDto `json:"rooms"`
}

type AvailableRoomDto struct {
	Adults   []AvailableRoomAdultDto    `json:"adults"`
	Children []AvailableRoomChildrenDto `json:"children"`
}

type AvailableRoomAdultDto struct {
	Title      string                    `json:"title"`
	FirstName  string                    `json:"firstName"`
	LastName   string                    `json:"lastName"`
	NationalId string                    `json:"nationalId,omitempty"`
	Cellphone  string                    `json:"cellphone,omitempty"`
	Passport   *AvailableRoomPassportDto `json:"passport,omitempty"`
}

type AvailableRoomChildrenDto struct {
	Title      string                    `json:"title"`
	FirstName  string                    `json:"firstName"`
	LastName   string                    `json:"lastName"`
	NationalId string                    `json:"nationalId,omitempty"`
	Passport   *AvailableRoomPassportDto `json:"passport,omitempty"`
}

type AvailableRoomPassportDto struct {
	Number           string `json:"number,omitempty"`
	ExpiryDate       string `json:"expiryDate,omitempty"`
	CountryResidency string `json:"countryResidency,omitempty"`
}

func (a AvailableDto) Validate() error {
	for _, room := range a.Rooms {
		for _, adult := range room.Adults {
			if adult.NationalId == "" && (adult.Passport == nil || adult.Passport.Number == "") {
				return common.NationalityUnknown
			}
		}
		for _, child := range room.Children {
			if child.NationalId == "" && (child.Passport == nil || child.Passport.Number == "") {
				return common.NationalityUnknown
			}
		}
	}
	return validation.ValidateStruct(&a,
		validation.Field(&a.HotelId, validation.Required),
		validation.Field(&a.SessionId, validation.Required),
		validation.Field(&a.PhoneNumber, validation.Required),
		validation.Field(&a.CheckOut, validation.Required),
		validation.Field(&a.CheckIn, validation.Required),
	)
}
