package dtos

type BaseResponse struct {
	Status string `json:"status"`
	Error  bool   `json:"error"`
}
