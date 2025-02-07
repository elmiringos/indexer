package rabbitmq

import (
	"encoding/json"

	"github.com/elmiringos/indexer/producer/pkg/logger"

	ampq "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Publisher struct {
	conn *ampq.Connection
	log  *zap.Logger
}

func NewPublisher(url string) *Publisher {
	log := logger.GetLogger()
	if log == nil {
		panic("logger is not initialized")
	}

	conn, err := ampq.Dial(url)
	if err != nil {
		log.Fatal("err in connecting to ampq", zap.Error(err), zap.String("URL", url))
	}

	return &Publisher{
		conn: conn,
		log:  log,
	}
}

func (p *Publisher) CreateChannel() *ampq.Channel {
	channel, err := p.conn.Channel()
	if err != nil {
		p.log.Fatal("err in creating channel to ampq", zap.Error(err))
	}

	return channel
}

func (p *Publisher) MakeNewQueueAndExchange(exchange ExchangeName, routingKey RoutingKey, queueType QueueType) (*ampq.Queue, error) {
	channel := p.CreateChannel()

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

func (p *Publisher) PublishMessage(channel *ampq.Channel, exchange ExchangeName, routingKey RoutingKey, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = channel.Publish(
		string(exchange),
		string(routingKey),
		false,
		false,
		ampq.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}

func (p *Publisher) CloseConnection() error {
	err := p.conn.Close()
	return err
}
