package database

import (
	"context"
	"fmt"
	"github.com/KonstantinPronin/email-sending-service/internal/sender/notification"
	"github.com/KonstantinPronin/email-sending-service/pkg/infrastructure"
	"github.com/KonstantinPronin/email-sending-service/pkg/model"
	"go.uber.org/zap"
)

type MongoDbClient struct {
	db     *infrastructure.Database
	logger *zap.Logger
}

func (m *MongoDbClient) Save(notif *model.Notification) error {
	_, err := m.db.GetConnection().Collection("notifications").InsertOne(context.Background(), notif)
	if err != nil {
		m.logger.Error(fmt.Sprintf("insert error: %s", err.Error()))
		return err
	}

	m.logger.Info(fmt.Sprintf("Message %s was saved to database", notif.ID))

	return nil
}

func NewMongoDbClient(
	db *infrastructure.Database,
	logger *zap.Logger) notification.Repository {
	return &MongoDbClient{
		db:     db,
		logger: logger,
	}
}
