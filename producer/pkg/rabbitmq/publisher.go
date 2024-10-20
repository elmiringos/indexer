package rabbitmq

import (
	"encoding/json"

	"github.com/elmiringos/indexer/producer/pkg/logger"

	"github.com/ethereum/go-ethereum/core/types"
	ampq "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Publisher struct {
	conn    *ampq.Connection
	channel *ampq.Channel
	log     *zap.Logger
}

func NewPublisher(url string) *Publisher {
	log := logger.GetLogger()
	if log == nil {
		panic("logger is not initialized")
	}

	conn, err := ampq.Dial(url)
	if err != nil {
		log.Fatal("err in connecting to ampq", zap.Error(err))
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal("err in creating channel to ampq", zap.Error(err))
	}

	return &Publisher{
		conn:    conn,
		channel: channel,
		log:     log,
	}
}

func (p *Publisher) MakeNewQueueAndExchange(exchange ExchangeName, routingKey RoutingKey, queueType QuqueType) (*ampq.Queue, error) {
	err := p.channel.ExchangeDeclare(
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

	queue, err := p.channel.QueueDeclare(
		string(queueType), // name
		false,             // durable
		false,             // auto-delete
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)

	err = p.channel.QueueBind(
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

func (p *Publisher) PublishBlockMessage(block *types.Block) error {
	body, err := json.Marshal(block.Header())
	if err != nil {
		return err
	}

	// Publish the message to the exchange
	err = p.channel.Publish(
		string(BlockExchange), // exchange
		string(BlockRoute),    // routing key
		false,                 // mandatory
		false,                 // immediate
		ampq.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}

func (p *Publisher) PublishTransactionMessage(transaction *types.Transaction) error {
	body, err := transaction.MarshalJSON()
	if err != nil {
		return err
	}

	// Publish the message to the exchange
	err = p.channel.Publish(
		string(TransactionExchange),
		string(TransactionRoute),
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
	err := p.channel.Close()
	return err
}

func (p *Publisher) CloseChannel() error {
	err := p.channel.Close()
	return err
}
