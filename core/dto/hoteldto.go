package dto

import (
	"hotel-engine/utils/indraframework"
	"time"
)

type HotelDto struct {
	PlaceID         string             `json:"placeId"`
	RoomID          string             `json:"roomId"`
	Type            string             `json:"type"`
	Kind            string             `json:"kind"`
	MinNight        int                `json:"min_night"`
	Code            int                `json:"code"`
	ReservationType string             `json:"reservation_type"`
	PaymentType     string             `json:"payment_type"`
	Name            string             `json:"name"`
	NameEn          string             `json:"nameEn"`
	Description     string             `json:"description"`
	Region          string             `json:"region"`
	SuitableFor     []string           `json:"suitable_for"`
	Images          []string           `json:"images"`
	Capacity        int                `json:"capacity"`
	CheckIn         time.Time          `json:"checkIn"`
	CheckOut        time.Time          `json:"checkOut"`
	CheckInTime     string             `json:"checkInTime"`
	CheckOutTime    string             `json:"checkOutTime"`
	Address         string             `json:"address"`
	Price           int64              `json:"price"`
	OldPrice        int64              `json:"oldPrice"`
	DiscountPercent int64              `json:"discountPercent"`
	DiscountPrice   int64              `json:"DiscountPrice"`
	City            string             `json:"city"`
	Province        string             `json:"province,omitempty"`
	CityEn          string             `json:"cityEn,omitempty"`
	ProvinceEn      string             `json:"provinceEn,omitempty"`
	GeoLocation     string             `json:"geoLocation,omitempty"`
	Tags            []string           `json:"tags"`
	Verified        bool               `json:"verified"`
	RateReviewScore float64            `json:"rateReviewScore"`
	Star            int                `json:"star"`
	RateReviewCount int                `json:"rateReviewCount"`
	Amenities       []HotelAmenityDto  `json:"amenities"`
	Places          []HotelPlaceDto    `json:"places"`
	Badges          []HotelBadgeDto    `json:"badges"`
	FAQ             HotelFAQDetailsDto `json:"faq"`
	CreatedAt       *time.Time         `json:"createdAt"`
	UpdatedAt       *time.Time         `json:"updatedAt"`

	UnavailableAmenities []HotelAmenityDto              `json:"unavailableAmenities"`
	Rooms                []RoomOptionDto                `json:"options,omitempty"`
	SessionId            string                         `json:"sessionId,omitempty"`
	Error                *indraframework.IndraException `json:"error"`
	Country              string                         `json:"country"`
	CountryEn            string                         `json:"countryEn"`
	CountryCode          string                         `json:"countryCode"`
	Sort                 float64                        `json:"sort"`
	Seo                  HotelSeoDto                    `json:"seo"`
}

type HotelSeoDto struct {
	Title           string `json:"title"`
	H1              string `json:"h1"`
	Description     string `json:"description"`
	Robots          string `json:"robots"`
	Canonical       string `json:"canonical"`
	MetaDescription string `json:"metaDescription"`
}

type HotelFAQDetailsDto struct {
	Title   string        `json:"title"`
	FAQList []HotelFAQDto `json:"faq"`
}

func (a *HotelDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}

func (a *HotelDto) SetProvince(province string) {
	a.Province = province
}

func (a *HotelDto) SetType(hotelType string) {
	a.Type = hotelType
}
