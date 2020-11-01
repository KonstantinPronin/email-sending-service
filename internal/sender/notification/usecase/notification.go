package usecase

import (
	"github.com/KonstantinPronin/email-sending-service/internal"
	"github.com/KonstantinPronin/email-sending-service/internal/sender/notification"
	"github.com/KonstantinPronin/email-sending-service/pkg/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type Notification struct {
	sender     internal.Transmitter
	repository notification.Repository
	logger     *zap.Logger
}

func (n *Notification) Send(notif *model.Notification) error {
	if notif.ID == "" {
		notif.ID = uuid.New().String()
	}

	if notif.CreatedAt == "" {
		notif.CreatedAt = time.Now().Format(time.RFC3339)
	}

	if !notif.SentStatus {
		if err := n.sender.Transfer(notif); err != nil {
			return err
		}
	}

	if err := n.repository.Save(notif); err != nil {
		return err
	}

	return nil
}

func NewNotification(
	sender internal.Transmitter,
	repository notification.Repository,
	logger *zap.Logger) notification.Usecase {
	return &Notification{
		sender:     sender,
		repository: repository,
		logger:     logger,
	}
}
