package dto

type HotelPlaceDto struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	GeoLocation string  `json:"geo_location"`
	Distance    float64 `json:"distance"`
	Priority    int     `json:"priority"`
}
