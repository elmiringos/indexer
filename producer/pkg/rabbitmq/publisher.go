package rabbitmq

import (
	"encoding/json"
	"fmt"

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

	conn, err := ampq.Dial(url)
	if err != nil {
		panic(fmt.Errorf("err in connecting to ampq: %v", err))
	}

	channel, err := conn.Channel()
	if err != nil {
		panic(fmt.Errorf("err in creating channel to ampq: %v", err))
	}

	return &Publisher{
		conn:    conn,
		channel: channel,
		log:     log,
	}
}

func (p *Publisher) MakeNewQueueAndExchange(exchange ExchangeName, queueType QuqueType) (*ampq.Queue, error) {
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

	return &queue, nil
}

func (p *Publisher) PublishBlockMessage(block *types.Block) error {
	body, err := json.Marshal(block)
	if err != nil {
		return err
	}

	// Publish the message to the exchange
	err = p.channel.Publish(
		string(BlockExchange), // exchange
		string(BlockQuque),    // routing key
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
	body, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	// Publish the message to the exchange
	err = p.channel.Publish(
		string(TransactionExchange),
		string(TransactionQuque),
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
