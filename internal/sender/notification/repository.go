package notification

import "github.com/KonstantinPronin/email-sending-service/pkg/model"

type Repository interface {
	Save(notif *model.Notification) error
}
