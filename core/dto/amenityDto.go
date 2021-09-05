package dto

import "hotel-engine/utils/indraframework"

type HotelAmenityDto struct {
	ID                int                            `json:"id"`
	Name              string                         `json:"name"`
	NameEn            string                         `json:"nameEn"`
	GroupId           int                            `json:"groupId"`
	IconUrl           string                         `json:"iconUrl"`
	AmenityCategoryID *uint                          `json:"amenityCategoryID"`
	AmenityCategory   *AmenityCategoryDto            `json:"amenityCategory"`
	Error             *indraframework.IndraException `json:"error"`
}

func (a *HotelAmenityDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
