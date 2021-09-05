package dtos

type HotelsListResponse struct {
	Status string           `json:"status"`
	Error  bool             `json:"error"`
	Result HotelsListResult `json:"result"`
}

type HotelsListResult struct {
	HotelsList []HotelItem `json:"hotelsList"`
	Total      int         `json:"totalCount"`
}

type HotelItem struct {
	Id                string            `json:"id"`
	GiataId           string            `json:"giataId"`
	Name              HotelItemName     `json:"name"`
	CityName          HotelItemCityName `json:"cityName"`
	MetaDescription   string            `json:"metaDescription"`
	SeoTitle          string            `json:"seoTitle"`
	AccommodationType string            `json:"accommodation"`
}

type HotelItemName struct {
	EN      string `json:"en"`
	ENIndex string `json:"enIndex"`
}

type HotelItemCityName struct {
	EN string `json:"en"`
}
