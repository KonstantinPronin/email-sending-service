package queue

import (
	"fmt"
	"github.com/KonstantinPronin/email-sending-service/internal/sender/notification"
	"github.com/KonstantinPronin/email-sending-service/pkg/constants"
	"github.com/KonstantinPronin/email-sending-service/pkg/infrastructure"
	"github.com/KonstantinPronin/email-sending-service/pkg/model"
	"github.com/mailru/easyjson"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Worker = func(delivery amqp.Delivery)

type RabbitMqClient struct {
	queue   *infrastructure.Queue
	pool    chan Worker
	usecase notification.Usecase
	logger  *zap.Logger
}

func NewRabbitMqClient(
	queue *infrastructure.Queue,
	usecase notification.Usecase,
	logger *zap.Logger) *RabbitMqClient {
	return &RabbitMqClient{
		queue:   queue,
		usecase: usecase,
		logger:  logger,
	}
}

func (r *RabbitMqClient) ListenAndServe() error {
	r.logger.Info("Starting queue listener")

	r.initPool()
	defer close(r.pool)

	for {
		msgs, err := r.queue.Consume(false, false, false, false, nil)
		if err != nil {
			r.logger.Error(err.Error())
			return err
		}

		for delivery := range msgs {
			worker := <-r.pool
			go worker(delivery)
		}
	}
}

func (r *RabbitMqClient) initPool() {
	r.pool = make(chan Worker, constants.MaxSendingWorkers)

	for i := 0; i < constants.MaxSendingWorkers; i++ {
		r.pool <- r.work
	}
}

func (r *RabbitMqClient) work(delivery amqp.Delivery) {
	notif := new(model.Notification)

	if err := easyjson.Unmarshal(delivery.Body, notif); err != nil {
		r.logger.Error(fmt.Sprintf("wrong message format: %s", err.Error()))

		if err = delivery.Reject(false); err != nil {
			r.logger.Error(fmt.Sprintf("ack error: %s", err.Error()))
		}
	}

	if err := r.usecase.Send(notif); err != nil {
		body, err := notif.MarshalJSON()
		if err == nil {
			delivery.Body = body
		}

		if err = delivery.Nack(false, true); err != nil {
			r.logger.Error(fmt.Sprintf("ack error: %s", err.Error()))
		}
	}

	if err := delivery.Ack(false); err != nil {
		r.logger.Error(fmt.Sprintf("ack error: %s", err.Error()))
	}

	r.pool <- r.work
}
