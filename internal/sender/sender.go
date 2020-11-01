package sender

import (
	"github.com/KonstantinPronin/email-sending-service/internal"
	mq "github.com/KonstantinPronin/email-sending-service/internal/sender/notification/delivery/queue"
	"github.com/KonstantinPronin/email-sending-service/internal/sender/notification/repository/database"
	"github.com/KonstantinPronin/email-sending-service/internal/sender/notification/repository/smtp"
	"github.com/KonstantinPronin/email-sending-service/internal/sender/notification/usecase"
	"github.com/KonstantinPronin/email-sending-service/pkg/infrastructure"
	"go.uber.org/zap"
)

type Sender struct {
	listener *mq.RabbitMqClient
}

func (s *Sender) Start() error {
	return s.listener.ListenAndServe()
}

func NewSender(
	db *infrastructure.Database,
	queue *infrastructure.Queue,
	smtpConf *infrastructure.SmtpConf,
	logger *zap.Logger) internal.Application {

	sender := smtp.NewEmail(smtpConf, logger)
	rep := database.NewMongoDbClient(db, logger)
	uc := usecase.NewNotification(sender, rep, logger)
	listener := mq.NewRabbitMqClient(queue, uc, logger)

	return &Sender{listener: listener}
}
