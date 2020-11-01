package acceptor

import (
	"github.com/KonstantinPronin/email-sending-service/internal"
	"github.com/KonstantinPronin/email-sending-service/internal/acceptor/notification/delivery/http"
	"github.com/KonstantinPronin/email-sending-service/internal/acceptor/notification/repository/database"
	mq "github.com/KonstantinPronin/email-sending-service/internal/acceptor/notification/repository/queue"
	"github.com/KonstantinPronin/email-sending-service/internal/acceptor/notification/usecase"
	"github.com/KonstantinPronin/email-sending-service/pkg/infrastructure"
	"github.com/labstack/echo"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
)

type Acceptor struct {
	e      *echo.Echo
	logger *zap.Logger
	port   string
}

func (a *Acceptor) Start() error {
	return a.e.Start(a.port)
}

func NewAcceptor(e *echo.Echo,
	db *infrastructure.Database,
	queue *infrastructure.Queue,
	logger *zap.Logger,
	port string) internal.Application {

	sanitizer := bluemonday.UGCPolicy()
	rep := database.NewMongoDbClient(db, logger)
	sender := mq.NewRabbitMqClient(queue, logger)
	uc := usecase.NewNotification(sender, rep, logger)

	http.NewAcceptorHandler(e, sanitizer, uc, logger)

	return &Acceptor{
		e:      e,
		logger: logger,
		port:   port,
	}
}
