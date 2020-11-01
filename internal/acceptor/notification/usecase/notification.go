package usecase

import (
	"fmt"
	"github.com/KonstantinPronin/email-sending-service/internal"
	"github.com/KonstantinPronin/email-sending-service/internal/acceptor/notification"
	"github.com/KonstantinPronin/email-sending-service/pkg/constants"
	"github.com/KonstantinPronin/email-sending-service/pkg/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type Notification struct {
	logger     *zap.Logger
	sender     internal.Transmitter
	repository notification.Repository
}

func (n *Notification) Get(id string) (*model.Notification, error) {
	return n.repository.Get(id)
}

func (n *Notification) GetList(page, perPage int64) (model.NotificationList, int64, error) {
	if page < 1 || perPage < 1 || perPage > constants.MaxPageSize {
		return nil, 0, model.NewInvalidArgument("wrong page params")
	}

	return n.repository.GetList(page, perPage)
}

func (n *Notification) Accept(notif *model.Notification) (string, error) {
	if len(notif.To) == 0 || len(notif.Message) == 0 {
		return "", model.NewInvalidArgument("to or message field is empty")
	}

	notif.ID = uuid.New().String()
	notif.CreatedAt = time.Now().Format(time.RFC3339)

	err := n.sender.Transfer(notif)
	if err != nil {
		return "", err
	}

	n.logger.Info(fmt.Sprintf("Message %s was accepted", notif.ID))
	return notif.ID, nil
}

func NewNotification(
	sender internal.Transmitter,
	repository notification.Repository,
	logger *zap.Logger) notification.Usecase {
	return &Notification{
		logger:     logger,
		sender:     sender,
		repository: repository,
	}
}
