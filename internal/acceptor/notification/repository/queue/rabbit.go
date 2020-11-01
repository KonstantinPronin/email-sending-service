package queue

import (
	"fmt"
	"github.com/KonstantinPronin/email-sending-service/internal"
	"github.com/KonstantinPronin/email-sending-service/pkg/constants"
	"github.com/KonstantinPronin/email-sending-service/pkg/infrastructure"
	"github.com/KonstantinPronin/email-sending-service/pkg/model"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMqClient struct {
	queue  *infrastructure.Queue
	logger *zap.Logger
}

func (r *RabbitMqClient) Transfer(notif *model.Notification) error {
	body, err := notif.MarshalJSON()
	if err != nil {
		r.logger.Error(fmt.Sprintf("marshalling error: %s", err.Error()))
		return err
	}

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  constants.Json,
		Body:         body,
	}

	err = r.queue.Publish(false, false, msg)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Queue connection error: %s", err.Error()))
		return err
	}

	return nil
}

func NewRabbitMqClient(
	queue *infrastructure.Queue,
	logger *zap.Logger) internal.Transmitter {
	return &RabbitMqClient{
		queue:  queue,
		logger: logger,
	}
}
