package internal

import (
	"github.com/KonstantinPronin/email-sending-service/pkg/model"
)

type Transmitter interface {
	Transfer(notif *model.Notification) error
}
