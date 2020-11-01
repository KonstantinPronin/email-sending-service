package notification

import (
	"github.com/KonstantinPronin/email-sending-service/pkg/model"
)

type Repository interface {
	Get(id string) (*model.Notification, error)
	GetList(page, perPage int64) (model.NotificationList, int64, error)
}
