package dto

import "hotel-engine/utils/indraframework"

type OrderEnquiryResponseDto struct {
	AllowedRefundPaymentMethods []string                       `json:"allowedRefundPaymentMethods"`
	Items                       []OrderEnquiryItemDto          `json:"items"`
	Error                       *indraframework.IndraException `json:"error"`
}

type OrderEnquiryItemDto struct {
	ProviderId          string                      `json:"providerId"`
	ProductProviderType string                      `json:"productProviderType"`
	Destination         string                      `json:"destination"`
	DestinationName     string                      `json:"destinationName"`
	Items               []OrderEnquiryItemOptionDto `json:"items"`
}

type OrderEnquiryItemOptionDto struct {
	ReferenceCode        string                                     `json:"referenceCode"`
	IsRefundable         bool                                       `json:"isRefundable"`
	PaidAmount           int64                                      `json:"paidAmount"`
	TotalPenaltyAmount   int64                                      `json:"totalPenaltyAmount"`
	RefundableAmount     int64                                      `json:"refundableAmount"`
	RefundableType       string                                     `json:"refundableType"`
	RefundableStatus     string                                     `json:"refundableStatus"`
	RefundStatus         string                                     `json:"refundStatus"`
	PassengerInformation OrderEnquiryItemOptionPassengerInformation `json:"passengerInformation"`
}

type OrderEnquiryItemOptionPassengerInformation struct {
	Title           string `json:"title"`
	Name            string `json:"name"`
	LastName        string `json:"lastName"`
	NamePersian     string `json:"namePersian"`
	LastNamePersian string `json:"lastNamePersian"`
}

func (a *OrderEnquiryResponseDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
