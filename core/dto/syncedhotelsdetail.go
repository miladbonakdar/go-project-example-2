package dto

import (
	"hotel-engine/utils/indraframework"
	"time"
)

type SyncedHotelsDetail struct {
	TotalSynced     int                            `json:"totalSynced"`
	TotalNotSynced  int                            `json:"totalNotSynced"`
	TotalHotels     int                            `json:"totalHotels"`
	SyncedHotels    []HotelSyncDetail              `json:"syncedHotels"`
	NotSyncedHotels []HotelSyncDetail              `json:"notSyncedHotels"`
	Error           *indraframework.IndraException `json:"error"`
}

func (a *SyncedHotelsDetail) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}

type HotelSyncDetail struct {
	PlaceID   string    `json:"placeId"`
	Price     int64     `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
