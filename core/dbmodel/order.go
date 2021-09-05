package dbmodel

import (
	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	ProviderOrderId        string      `gorm:"column:ProviderOrderId;type:nvarchar(50);not null"`
	IndraOrderId           int64       `gorm:"column:IndraOrderId;not null"`
	HotelID                uint        `gorm:"column:HotelID;not null"`
	ProviderHotelId        string      `gorm:"column:ProviderHotelId;type:nvarchar(50);not null"`
	NonRefundable          bool        `gorm:"column:NonRefundable;not null"`
	GeneralPolicies        string      `gorm:"column:GeneralPolicies;not null,type:nvarchar(2500)"`
	TotalPrice             int64       `gorm:"column:TotalPrice;not null"`
	Provider               string      `gorm:"column:Provider;type:nvarchar(50)"`
	ProviderName           string      `gorm:"column:ProviderName;type:nvarchar(50)"`
	Currency               string      `gorm:"column:Currency;type:nvarchar(50);not null"`
	MealPlan               string      `gorm:"column:MealPlan;type:nvarchar(100);not null"`
	RestrictedMarkupAmount int64       `gorm:"column:RestrictedMarkupAmount;not null"`
	RestrictedMarkupType   string      `gorm:"column:RestrictedMarkupType;type:nvarchar(100);not null"`
	Status                 string      `gorm:"column:Status;type:nvarchar(50);not null"`
	Rooms                  []OrderRoom `gorm:"foreignKey:OrderID"`
	TransactionStatus      string      `gorm:"column:TransactionStatus;type:nvarchar(50);null"`
	TransactionRequestId   string      `gorm:"column:TransactionRequestId;type:nvarchar(50);null"`
	TransactionIds         string      `gorm:"column:TransactionIds;type:nvarchar(2500);null"`
	Confirmed              bool        `gorm:"column:Confirmed;not null;default:0"`

	RefundRequestId          int64   `gorm:"column:RefundRequestId;not null;default:0"`
	ApplicantRefundRequestId int64   `gorm:"column:ApplicantRefundRequestId;not null;default:0"`
	ApplicantOrderId         int64   `gorm:"column:ApplicantOrderId;not null;default:0"`
	PaidAmount               float32 `gorm:"column:PaidAmount;not null;type:decimal(10,2);default:0.0"`
	ReferenceCode            string  `gorm:"column:ReferenceCode;not null;type:nvarchar(50)"`
	RefundStatus             string  `gorm:"column:RefundStatus;not null;type:nvarchar(50)"`
	RefundableAmount         float32 `gorm:"column:RefundableAmount;not null;type:decimal(10,2);default:0.0"`
	TotalPenaltyAmount       float32 `gorm:"column:TotalPenaltyAmount;not null;type:decimal(10,2);default:0.0"`
}

func (h *Order) UpdateStatus(status string) {
	h.Status = status
}

func (h *Order) UpdateRefundRequestId(refundRequestId, applicantRefundRequestId, applicantOrderId int64) {
	h.RefundRequestId = refundRequestId
	h.ApplicantOrderId = applicantOrderId
	h.ApplicantRefundRequestId = applicantRefundRequestId
}

func (h *Order) UpdateRefundResult(paidAmount float32, referenceCode string,
	refundStatus string,
	refundableAmount float32, totalPenaltyAmount float32) {
	h.PaidAmount = paidAmount
	h.ReferenceCode = referenceCode
	h.RefundStatus = refundStatus
	h.RefundableAmount = refundableAmount
	h.TotalPenaltyAmount = totalPenaltyAmount
}

func (h *Order) UpdateConfirmed(confirmed bool) {
	h.Confirmed = confirmed
}

func (h *Order) UpdateTransaction(transactionStatus, transactionRequestId, transactionIds string) {
	h.TransactionStatus = transactionStatus
	h.TransactionRequestId = transactionRequestId
	h.TransactionIds = transactionIds
}
