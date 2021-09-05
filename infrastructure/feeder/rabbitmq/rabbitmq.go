package rabbitmq

import (
	"encoding/json"
	"fmt"
	"hotel-engine/core"
	"hotel-engine/core/dto"
	"hotel-engine/core/messaging"
	"hotel-engine/infrastructure/logger"
	"strconv"

	"github.com/streadway/amqp"
)

type rabbitmq struct {
	feed      string
	seed      string
	alias     string
	messaging messaging.Bus
}

//NewFeeder ...
func NewFeeder(messaging messaging.Bus, feed string, seed string, alias string) (core.Feeder, error) {
	return &rabbitmq{
		messaging: messaging,
		feed:      feed,
		seed:      seed,
		alias:     alias,
	}, nil
}

//Feed ...
func (m *rabbitmq) Feed(hotels dto.ElasticUpdateRequest) error {
	data, _ := json.Marshal(hotels)
	logger.WithData(map[string]string{
		"hotels": strconv.Itoa(len(hotels.Places)),
	}).WithDevMessage("feeder -> Feed").Info("Feeder : Sent Feed Data ...")
	return m.messaging.PublishOnQueue(data, m.feed)
}

//Seed ...
func (m *rabbitmq) Seed(handle func(index string)) error {
	return m.messaging.Subscribe(m.seed, "fanout", fmt.Sprintf("%s_%s", m.seed, "Hotel"), "Hotel", func(d amqp.Delivery) {
		logger.WithDevMessage("feeder -> Seed").Info("Feeder : Start Seed Data ...")
		handle(string(d.Body[:]))
	})
}

//Alias ...
func (m *rabbitmq) Alias(index string) error {
	logger.Print("Feeder : End Feed Data ...")
	body, _ := json.Marshal(map[string]string{
		"index":   index,
		"service": "hotel",
	})
	return m.messaging.PublishOnQueue(body, m.alias)
}

//Close ...
func (m *rabbitmq) Close() error {
	if m.messaging != nil {
		m.messaging.Close()
	}
	logger.Info("rabbit connection closed")
	return nil
}
