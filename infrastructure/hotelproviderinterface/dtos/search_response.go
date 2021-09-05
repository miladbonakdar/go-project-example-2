package dtos

type SearchResponse struct {
	BaseResponse
	Result SearchResult `json:"result"`
}

type SearchResult struct {
	SessionId string `json:"sessionId"`
}
