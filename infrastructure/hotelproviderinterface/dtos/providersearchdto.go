package dtos

import "encoding/json"

type ProviderSearchDto struct {
	RequestSession      ProviderRequestSessionDto
	RequestSearchHotels ProviderSearchHotelsRequestDto
}

type ProviderRequestSessionDto struct {
	CheckIn     string                               `json:"checkIn"`
	CheckOut    string                               `json:"checkOut"`
	Rooms       []ProviderRequestSessionRoomDto      `json:"rooms"`
	Destination ProviderRequestSessionDestinationDto `json:"destination"`
}

func (s ProviderRequestSessionDto) ToJson() []byte {
	jsonValue, _ := json.Marshal(s)
	return jsonValue
}

type ProviderRequestSessionRoomDto struct {
	Adults   []int `json:"adults"`
	Children []int `json:"children"`
}
type ProviderRequestSessionDestinationDto struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type ProviderSearchHotelsRequestDto struct {
	SessionId string                      `json:"sessionId"`
	Limit     int64                       `json:"limit"`
	Skip      int64                       `json:"skip"`
	Sort      ProviderSearchHotelsSortDto `json:"sort"`
	Filters   []interface{}               `json:"filter"`
}

func (s ProviderSearchHotelsRequestDto) ToJson() []byte {
	jsonValue, _ := json.Marshal(s)
	return jsonValue
}

type ProviderSearchHotelsSortDto struct {
	Field string `json:"field"`
	Order int    `json:"order"`
}

type SearchHotelsStringFilter struct {
	Field string   `json:"field"`
	Value []string `json:"value"`
}

type SearchHotelsIntFilter struct {
	Field string `json:"field"`
	Value []int  `json:"value"`
}

type SearchHotelsFloatFilter struct {
	Field string    `json:"field"`
	Value []float32 `json:"value"`
}

type SearchHotelsIntRangeFilter struct {
	Field string    `json:"field"`
	Value [][]int64 `json:"value"`
}
