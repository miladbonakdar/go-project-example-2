package logic

import (
	"encoding/json"
	"fmt"
	"hotel-engine/core"
	"hotel-engine/core/common/logtags"
	"hotel-engine/core/dto"
	"hotel-engine/core/messaging"
	"hotel-engine/infrastructure/logger"
)

type orderEventDispatcher struct {
	client       messaging.Bus
	messageTopic string
}

type OrderRefundRequestFinalizedEvent struct {
	PenaltyAmount   float32                           `json:"PenaltyAmount"`
	OrderID         int64                             `json:"OrderId"`
	RefundRequestID int64                             `json:"RefundRequestId"`
	ItemDetails     []OrderRefundRequestFinalizedItem `json:"ItemDetails"`
}

type OrderRefundRequestFinalizedItem struct {
	IsSuccessful     bool    `json:"IsSuccessful"`
	Item             string  `json:"Item"`
	PenaltyAmount    float32 `json:"PenaltyAmount"`
	AlibabaRefundFee float32 `json:"AlibabaRefundFee"`
}

func (d *orderEventDispatcher) OrderRefundRequestFinalized(event dto.OrderRefundRequestFinalizedDto) {
	data := OrderRefundRequestFinalizedEvent{
		PenaltyAmount:   event.TotalPenaltyAmount,
		OrderID:         event.ApplicantOrderId,
		RefundRequestID: event.ApplicantRefundRequestId,
		ItemDetails: []OrderRefundRequestFinalizedItem{
			{
				IsSuccessful:     true,
				Item:             event.ProviderOrderId,
				PenaltyAmount:    event.TotalPenaltyAmount,
				AlibabaRefundFee: 0,
			},
		},
	}
	body, err := json.Marshal(data)

	if err != nil {
		logger.WithName(logtags.CannotCreateRefundEventError).ErrorException(err, "Cannot create order refund event object")
		return
	}
	exchangeType := "topic"
	err = d.client.Publish(body, d.messageTopic, exchangeType, d.messageTopic)
	if err != nil {
		logger.WithName(logtags.CannotPublishRefundEventError).ErrorException(err, "Cannot publish order refund event object")
	}
	logger.WithData(data).WithName(logtags.RefundRequestCompleted).
		Info(fmt.Sprintf("Refund request for order with id %s completed automatically", event.ProviderOrderId))
}

func NewOrderEventDispatcher(client messaging.Bus, topic string) core.OrderEventDispatcher {
	return &orderEventDispatcher{
		client:       client,
		messageTopic: topic,
	}
}
