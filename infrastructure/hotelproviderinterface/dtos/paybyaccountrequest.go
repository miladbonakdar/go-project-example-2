package dtos

import "encoding/json"

type PayByAccountRequest struct {
	RedirectUrl string `json:"redirectUrl"`
}

func NewPayByAccountRequest(url string) PayByAccountRequest {
	return PayByAccountRequest{
		RedirectUrl: url,
	}
}

func (s PayByAccountRequest) ToJson() []byte {
	jsonValue, _ := json.Marshal(s)
	return jsonValue
}
