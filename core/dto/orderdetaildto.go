package dto

import "hotel-engine/utils/indraframework"

type OrderDetailDto struct {
	OrderId                  uint                           `json:"OrderId"`
	ProviderOrderId          string                         `json:"ProviderOrderId"`
	IndraOrderId             int64                          `json:"IndraOrderId"`
	HotelID                  uint                           `json:"HotelID"`
	ProviderHotelId          string                         `json:"ProviderHotelId"`
	NonRefundable            bool                           `json:"NonRefundable"`
	GeneralPolicies          []string                       `json:"GeneralPolicies"`
	TotalPrice               int64                          `json:"TotalPrice"`
	Provider                 string                         `json:"Provider"`
	ProviderName             string                         `json:"ProviderName"`
	Currency                 string                         `json:"Currency"`
	MealPlan                 string                         `json:"MealPlan"`
	RestrictedMarkupAmount   int64                          `json:"RestrictedMarkupAmount"`
	RestrictedMarkupType     string                         `json:"RestrictedMarkupType"`
	Status                   string                         `json:"Status"`
	Rooms                    []OrderDetailRoomDto           `json:"Rooms"`
	Hotel                    HotelDto                       `json:"Hotel"`
	ApplicantRefundRequestId int64                          `json:"ApplicantRefundRequestId"`
	ApplicantOrderId         int64                          `json:"ApplicantOrderId"`
	PaidAmount               float32                        `json:"PaidAmount"`
	ReferenceCode            string                         `json:"ReferenceCode"`
	RefundStatus             string                         `json:"RefundStatus"`
	RefundableAmount         float32                        `json:"RefundableAmount"`
	TotalPenaltyAmount       float32                        `json:"TotalPenaltyAmount"`
	RefundRequestId          int64                          `json:"RefundRequestId"`
	Confirmed                bool                           `json:"Confirmed"`
	Error                    *indraframework.IndraException `json:"error"`
}

func (a *OrderDetailDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}

type OrderDetailRoomDto struct {
	OrderRoomId   uint   `json:"OrderRoomId"`
	OrderId       uint   `json:"OrderId"`
	PricePerNight int64  `json:"PricePerNight"`
	Price         int64  `json:"Price"`
	Name          string `json:"Name"`
	NameEn        string `json:"NameEn"`
}
