package infrastructure

import (
	"fmt"
	"github.com/KonstantinPronin/email-sending-service/pkg/constants"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"time"
)

type Queue struct {
	ch   *amqp.Channel
	url  string
	dest string
}

func (q *Queue) Publish(mandatory, immediate bool, msg amqp.Publishing) error {
	err := q.ch.Publish("", q.dest, mandatory, immediate, msg)

	if err == nil {
		return nil
	}

	err = q.Reconnect()
	if err != nil {
		return err
	}

	return q.ch.Publish("", q.dest, mandatory, immediate, msg)
}

func (q *Queue) Consume(autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	delivery, err := q.ch.Consume(q.dest, "", autoAck, exclusive, noLocal, noWait, args)

	if err == nil {
		return delivery, nil
	}

	err = q.Reconnect()
	if err != nil {
		return nil, err
	}

	return q.ch.Consume(q.dest, "", autoAck, exclusive, noLocal, noWait, args)
}

func (q *Queue) Connect() error {
	conn, err := amqp.Dial(q.url)
	if err != nil {
		return err
	}

	q.ch, err = conn.Channel()
	if err != nil {
		return err
	}

	return nil
}

func (q *Queue) Close() error {
	if q.ch != nil {
		return q.ch.Close()
	}

	return nil
}

func (q *Queue) Reconnect() error {
	for attempt := constants.MaxConnectAttempts; attempt > 0; attempt-- {
		err := q.Connect()

		if err == nil {
			return nil
		}

		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("queue connection timeout")
}

func InitQueue(path string) (*Queue, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	dest := viper.GetString("queue")
	url := viper.GetString("url")
	queue := &Queue{
		ch:   nil,
		url:  url,
		dest: dest,
	}

	if err := queue.Reconnect(); err != nil {
		return nil, err
	}

	if _, err := queue.ch.QueueDeclare(dest, true, false, false, false, nil); err != nil {
		return nil, err
	}

	return queue, nil
}
