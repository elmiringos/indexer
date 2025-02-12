package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/elmiringos/indexer/indexer-core/pkg/logger"
)

type Consumer struct {
	conn *amqp091.Connection
	log  *zap.Logger
}

func NewConsumer(url string) *Consumer {
	log := logger.GetLogger()
	if log == nil {
		panic("logger is not initialized")
	}

	conn, err := amqp091.Dial(url)
	if err != nil {
		log.Fatal("err in connecting to ampq", zap.Error(err), zap.String("URL", url))
	}

	return &Consumer{
		conn: conn,
		log:  log,
	}
}

func (c *Consumer) CreateChannel() *amqp091.Channel {
	channel, err := c.conn.Channel()
	if err != nil {
		c.log.Fatal("err in creating channel to ampq", zap.Error(err))
	}

	return channel
}

func (c *Consumer) MakeNewQueueAndExchange(queueType QueueType, exchange ExchangeName, routingKey RoutingKey) (*amqp091.Queue, error) {
	channel := c.CreateChannel()

	err := channel.ExchangeDeclare(
		string(exchange), // name
		"direct",         // type
		false,            // not durable
		false,            // auto-delete
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		string(queueType), // name
		false,             // durable
		false,             // auto-delete
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return nil, err
	}

	err = channel.QueueBind(
		queue.Name,         // name of the queue
		string(routingKey), // routing key (messages with this key will be routed to this queue)
		string(exchange),   // name of the exchange
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return nil, err
	}

	return &queue, nil
}

func (c *Consumer) Consume(channel *amqp091.Channel, queue QueueType) <-chan amqp091.Delivery {

	msgs, err := channel.Consume(
		string(queue),
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		c.log.Fatal("Error in cosuming messages from rabbitmq", zap.Error(err), zap.String("queue", string(queue)))
	}

	return msgs
}

func (c *Consumer) CloseConnection() error {
	return c.conn.Close()
}
