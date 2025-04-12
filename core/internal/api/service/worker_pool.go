package service

import (
	"context"
	"sync"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// MessageProcessor is an interface for processing messages
type MessageProcessor interface {
	Process(context.Context, []byte) error
}

// WorkerPool represents a generic worker pool for processing messages
type WorkerPool struct {
	processor   MessageProcessor
	log         *zap.Logger
	workerCount int
}

// NewWorkerPool creates a new WorkerPool
func NewWorkerPool(processor MessageProcessor, log *zap.Logger, workerCount int) *WorkerPool {
	return &WorkerPool{
		processor:   processor,
		log:         log,
		workerCount: workerCount,
	}
}

// Start starts the worker pool
func (p *WorkerPool) Start(msgs <-chan amqp091.Delivery) {
	p.log.Info("Starting worker pool", zap.Int("worker_count", p.workerCount))

	var wg sync.WaitGroup

	for i := 0; i < p.workerCount; i++ {
		wg.Add(1)
		go p.worker(i, &wg, msgs)
	}

	wg.Wait()
}

// worker is a worker function that processes messages
func (p *WorkerPool) worker(id int, wg *sync.WaitGroup, msgs <-chan amqp091.Delivery) {
	defer wg.Done()

	for msg := range msgs {
		p.log.Info("Worker processing message", zap.Int("worker_id", id))

		err := p.processor.Process(context.Background(), msg.Body)
		if err != nil {
			p.log.Error(
				"Failed to process message",
				zap.Error(err),
				zap.Int("worker_id", id),
			)

			// Send nack with requeue=true to retry processing the message
			if nackErr := msg.Nack(false, true); nackErr != nil {
				p.log.Fatal("Error sending nack message", zap.Error(nackErr))
			}
			continue
		}

		// Acknowledge the message after successful processing
		if ackErr := msg.Ack(false); ackErr != nil {
			p.log.Error("Failed to acknowledge message",
				zap.Error(ackErr),
				zap.Int("worker_id", id))
			continue
		}

		p.log.Info("Message processed successfully", zap.Int("worker_id", id), zap.String("msg_id", string(msg.MessageId)))
	}
}
