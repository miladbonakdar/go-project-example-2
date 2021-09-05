package dtos

import "encoding/json"

type RefundRequestDto struct {
	RefundRequestType   string          `json:"refundRequestType"`
	RefundPaymentMethod string          `json:"refundPaymentMethod"`
	OrderId             string          `json:"orderId"`
	Items               []RefundItemDto `json:"items"`
}

type RefundItemDto struct {
	ReferenceCode string `json:"referenceCode"`
}

func NewRefundRequestDto(orderId, referenceCode string) RefundRequestDto {
	return RefundRequestDto{
		RefundRequestType:   "Personal",
		RefundPaymentMethod: "UserAccount",
		OrderId:             orderId,
		Items: []RefundItemDto{
			{ReferenceCode: referenceCode},
		},
	}
}

func (s RefundRequestDto) ToJson() []byte {
	jsonValue, _ := json.Marshal(s)
	return jsonValue
}
