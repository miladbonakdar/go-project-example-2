package dtos

type SearchDirectResponse struct {
	BaseResponse
	Result SearchDirectResult `json:"result"`
}

type SearchDirectResult struct {
	SearchResult
	HotelId string `json:"hotelId"`
}
