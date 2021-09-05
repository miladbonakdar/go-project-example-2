package dtos

import (
	"encoding/json"
)

func (s ResultRequest) ToJson() []byte {
	jsonValue, _ := json.Marshal(s)
	return jsonValue
}

type ResultRequest struct {
	SessionId string `json:"sessionId"`
	Limit     int    `json:"limit"`
	Skip      int    `json:"skip"`
}
