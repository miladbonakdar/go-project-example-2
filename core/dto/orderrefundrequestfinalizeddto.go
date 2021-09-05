package dto

type OrderRefundRequestFinalizedDto struct {
	ApplicantRefundRequestId int64
	ApplicantOrderId         int64
	ProviderOrderId          string
	PaidAmount               float32
	RefundStatus             string
	RefundableAmount         float32
	TotalPenaltyAmount       float32
}
