package rabbitmq

import (
	ampq "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/elmiringos/indexer/indexer-core/pkg/logger"
)

type Subscriber struct {
	conn *ampq.Connection
	log  *zap.Logger
}

func NewSubscriber(url string) *Subscriber {
	log := logger.GetLogger()
	if log == nil {
		panic("logger is not initialized")
	}

	conn, err := ampq.Dial(url)
	if err != nil {
		log.Fatal("err in connecting to ampq", zap.Error(err), zap.String("URL", url))
	}

	return &Subscriber{
		conn: conn,
		log:  log,
	}
}
