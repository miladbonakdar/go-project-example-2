package logic

import (
	"encoding/json"
	"hotel-engine/core"
	"hotel-engine/core/common/logtags"
	"hotel-engine/core/dto"
	"hotel-engine/core/messaging"
	"hotel-engine/infrastructure/logger"
	"strings"

	"github.com/streadway/amqp"
)

var rateHandler RateReviewEventHandler

type RateReviewEventHandler struct {
	client       messaging.Bus
	hotelService core.HotelService
}

const hotelProductType = 2

func (r *RateReviewEventHandler) handleMessage(message dto.RateReviewEventDto) {
	if message.ProductType != hotelProductType {
		return
	}
	_ = r.hotelService.UpdateHotelRateReview(message)
}

func (r *RateReviewEventHandler) subscribeCallback(d amqp.Delivery) {
	var message dto.RateReviewEventDto
	err := json.Unmarshal(d.Body, &message)
	if err != nil {
		logger.WithName(logtags.CastRateReviewEventObjectError).
			ErrorException(err, "cannot cast body of the event to the slice of events")
	}
	r.handleMessage(message)
}

func NewRateReviewEventHandler(client messaging.Bus, service core.HotelService, subscribeString string) *RateReviewEventHandler {
	items := strings.Split(subscribeString, ",")
	rateHandler := &RateReviewEventHandler{
		client:       client,
		hotelService: service,
	}
	err := rateHandler.client.Subscribe(items[0], items[1], items[2], items[3], rateHandler.subscribeCallback)
	if err != nil {
		logger.WithName(logtags.CannotSubscribeToRateReviewQueue).
			FatalException(err, "error while trying to subscribe to rate and review queue")
	}
	return rateHandler
}
