package rabbitmq

import (
	ampq "github.com/rabbitmq/amqp091-go"
)

type BlockPublisher struct {
	conn     *ampq.Connection
	channel  *ampq.Channel
	exchange string
}

func NewBlockPublisher(url string) *BlockPublisher {
	conn, err := ampq.Dial(url)
	if err != nil {
		panic("err in connecting to ampq")
	}

	return &BlockPublisher{
		conn: conn,
	}
}
