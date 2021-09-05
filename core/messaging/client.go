package messaging

import (
	"fmt"
	"hotel-engine/core/common/logtags"
	"hotel-engine/infrastructure/logger"
	"time"

	"github.com/streadway/amqp"
)

type Bus interface {
	ConnectToBroker(connectionString string)
	Publish(msg []byte, exchangeName string, exchangeType string, messageType string) error
	PublishOnQueue(msg []byte, queueName string) error
	Subscribe(exchangeName string, exchangeType string, queueName string, consumerName string, worker func(amqp.Delivery)) error
	SubscribeToQueue(queueName string, consumerName string, consumerCount int, consumerSize int, consumer func(amqp.Delivery)) error
	Close()
}

type client struct {
	conn             *amqp.Connection
	rabbitCloseError chan *amqp.Error
	channels         []*channel
}

type channel struct {
	channel  *amqp.Channel
	exchange string
	queue    string
	key      string
	unbind   bool
}

// re-establish the connection to RabbitMQ in case
// the connection has died
//
func (c *client) rabbitConnector(uri string) {
	var rabbitErr *amqp.Error
	c.conn = c.connectToRabbitMQ(uri)
	for {
		rabbitErr = <-c.rabbitCloseError
		if rabbitErr != nil {
			c.conn = c.connectToRabbitMQ(uri)
			c.rabbitCloseError = make(chan *amqp.Error)
			c.conn.NotifyClose(c.rabbitCloseError)
		}
	}
}
func (c *client) connectToRabbitMQ(uri string) *amqp.Connection {
	for {
		conn, err := amqp.Dial(uri)

		if err == nil {
			return conn
		}
		logger.WithName(logtags.RabbitConnectionError).WithException(err).Error(err.Error())
		time.Sleep(3 * time.Second)
	}
}
func (c *client) ConnectToBroker(connectionString string) {
	c.rabbitCloseError = make(chan *amqp.Error)
	go c.rabbitConnector(connectionString)
	c.rabbitCloseError <- amqp.ErrClosed
}

//Publish ...
func (c *client) Publish(body []byte, exchangeName string, exchangeType string, messageType string) error {
	if c.conn == nil {
		panic("Tried to send message before connection was initialized. Don't do that.")
	}

	ch, err := c.conn.Channel()
	defer ch.Close()
	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Failed to register an Exchange")
	err = ch.Publish( // Publishes a message onto the queue.
		exchangeName,
		exchangeName,
		false,
		false,
		amqp.Publishing{
			Body: body, // Our JSON body as []byte
			Type: messageType,
		})
	//fmt.Printf("A message was sent: %v", body)
	return err
}

//PublishOnQueue ...
func (c *client) PublishOnQueue(body []byte, queueName string) error {
	if c.conn == nil {
		panic("Tried to send message before connection was initialized. Don't do that.")
	}
	ch, err := c.conn.Channel()
	defer ch.Close()

	queue, err := ch.QueueDeclare( // Declare a queue that will be created if not exists with some args
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	// Publishes a message onto the queue.
	err = ch.Publish(
		"",
		queue.Name,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	// fmt.Printf("A message was sent to queue %v: %v", queueName, body)
	return err
}

//Subscribe ...
func (c *client) Subscribe(exchangeName string, exchangeType string, queueName string, consumerName string, worker func(amqp.Delivery)) error {
	ch, err := c.conn.Channel()
	handleError(err, "Failed to open a channel")
	// defer ch.Close()

	c.channels = append(c.channels, &channel{
		channel:  ch,
		exchange: exchangeName,
		queue:    queueName,
		key:      exchangeName,
		unbind:   false,
	})

	var args map[string]interface{}

	if exchangeType == "x-delayed-message" {
		args = map[string]interface{}{
			"x-delayed-type": "topic",
		}
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		args,
	)

	handleError(err, "Failed to register an Exchange")

	queue, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Failed to register an Queue")

	err = ch.QueueBind(
		queue.Name,
		exchangeName,
		exchangeName,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("queue Bind: %s", err)
	}

	go consumeLoop(ch, worker, queue.Name, consumerName, 1, 1)

	return nil
}

func (c *client) SubscribeToQueue(queueName string, consumerName string, consumerCount int, consumerSize int, consumer func(amqp.Delivery)) error {
	ch, err := c.conn.Channel()
	handleError(err, "Failed to open a channel")

	c.channels = append(c.channels, &channel{
		channel: ch,
		unbind:  false,
	})

	queue, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Failed to register an Queue")

	if consumerCount > 1 {
		if err := ch.Qos(consumerSize, 0, false); err != nil {
			return err
		}

		for i := 1; i <= consumerCount; i++ {
			id := i
			go consumeLoop(ch, consumer, queue.Name, consumerName, consumerCount, id)
		}

	} else {
		go consumeLoop(ch, consumer, queue.Name, consumerName, 1, 1)
	}

	return nil
}

func (c *channel) close() {
	//Comment On LocalTest
	// if c.unbind {
	// 	c.channel.QueueUnbind(c.queue, c.exchange, c.key, nil)
	// }
	c.channel.Close()
}

//Close ...
func (c *client) Close() {
	if c.channels != nil {
		for _, c := range c.channels {
			c.close()
		}
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func consumeLoop(
	ch *amqp.Channel,
	handlerFunc func(d amqp.Delivery),
	qName string,
	consumerName string,
	consumerCount int,
	id int) {

	MSGs, err := ch.Consume(
		qName,
		fmt.Sprintf("%s (%d/%d)", consumerName, id, consumerCount),
		true,
		false,
		false,
		false,
		nil,
	)

	handleError(err, fmt.Sprintf("Failed to register a consumer (%d/%d)", id, consumerCount))

	for msg := range MSGs {
		// Invoke the handlerFunc func we passed as parameter.
		handlerFunc(msg)

		//TODO:: don't uncomment these lines.
		//if err := msg.Ack(false); err != nil {
		//	log.Println("unable to acknowledge the message, dropped", err)
		//}
	}

	logger.Info(fmt.Sprintf("[%d] Exiting ...", id))
}

func handleError(err error, msg string) {
	if err != nil {
		//fmt.Printf("%s: %s", msg, err)
		logger.WithName(logtags.RabbitUnknownError).Fatal(fmt.Sprintf("%s: %s", msg, err))
	}
}

func NewBusClient(connectionString string) Bus {
	c := &client{}
	c.ConnectToBroker(connectionString)
	return c
}
