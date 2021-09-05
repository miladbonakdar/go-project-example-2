package dto

type OrdersRefundStatusResponseDto struct {
	Success bool                                `json:"success"`
	Error   interface{}                         `json:"error"`
	Result  OrdersRefundStatusResponseDtoResult `json:"result"`
}

type OrdersRefundStatusResponseDtoResult struct {
	PageNumber int                                       `json:"pageNumber"`
	PageSize   int                                       `json:"pageSize"`
	TotalCount int                                       `json:"totalCount"`
	Items      []OrdersRefundStatusResponseDtoResultItem `json:"items"`
}

type OrdersRefundStatusResponseDtoResultItem struct {
	RefundRequestId      int64                                         `json:"refundRequestId"`
	OrderId              int64                                         `json:"orderId"`
	UserUniqueNumber     int64                                         `json:"userUniqueNumber"`
	RequestedRefundType  string                                        `json:"requestedRefundType"`
	RefundStatus         string                                        `json:"refundStatus"`
	RefundPaymentMethod  string                                        `json:"refundPaymentMethod"`
	CreationMethod       string                                        `json:"creationMethod"`
	TransactionRequestId int64                                         `json:"transactionRequestId"`
	CreatorUserId        int64                                         `json:"creatorUserId"`
	ReferenceCodes       string                                        `json:"referenceCodes"`
	Items                []OrdersRefundStatusResponseDtoResultItemItem `json:"items"`
}

type OrdersRefundStatusResponseDtoResultItemItem struct {
	Id                  int64   `json:"id"`
	ReferenceCode       string  `json:"referenceCode"`
	OrderLineItemId     int64   `json:"orderLineItemId"`
	ProductProviderType string  `json:"productProviderType"`
	ProviderPenalty     float32 `json:"providerPenalty"`
	AcceptedPenalty     float32 `json:"acceptedPenalty"`
	AlibabaRefundFee    float32 `json:"alibabaRefundFee"`
	TotalPenaltyAmount  float32 `json:"totalPenaltyAmount"`
	TotalAmount         float32 `json:"totalAmount"`
	PaidAmount          float32 `json:"paidAmount"`
	RefundableAmount    float32 `json:"refundableAmount"`
	IsRefunded          bool    `json:"isRefunded"`
}
