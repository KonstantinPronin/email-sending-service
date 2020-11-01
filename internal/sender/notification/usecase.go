package notification

import "github.com/KonstantinPronin/email-sending-service/pkg/model"

type Usecase interface {
	Send(notif *model.Notification) error
}
