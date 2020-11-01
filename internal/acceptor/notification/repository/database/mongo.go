package database

import (
	"context"
	"fmt"
	"github.com/KonstantinPronin/email-sending-service/internal/acceptor/notification"
	"github.com/KonstantinPronin/email-sending-service/pkg/infrastructure"
	"github.com/KonstantinPronin/email-sending-service/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"math"
)

type MongoDbClient struct {
	db     *infrastructure.Database
	logger *zap.Logger
}

func (m *MongoDbClient) Get(id string) (*model.Notification, error) {
	notif := new(model.Notification)
	filter := bson.M{"id": id}

	result := m.db.GetConnection().
		Collection("notifications").FindOne(context.Background(), filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, model.NewNotFoundError(result.Err().Error())
		}
		return nil, result.Err()
	}

	if err := result.Decode(notif); err != nil {
		m.logger.Error(fmt.Sprintf("Decode error: %s", err.Error()))
		return nil, err
	}

	return notif, nil
}

func (m *MongoDbClient) GetList(page, perPage int64) (model.NotificationList, int64, error) {
	result := model.NotificationList{}
	page = page - 1
	opt := options.Find()
	filter := bson.M{}

	total, err := m.db.GetConnection().
		Collection("notifications").CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	if page > int64(math.Ceil(float64(total)/float64(perPage))) {
		return nil, 0, model.NewInvalidArgument("page over limit")
	}

	opt.SetLimit(perPage)
	opt.SetSort(bson.M{"date": 1})
	opt.SetSkip((page) * perPage)

	cursor, err := m.db.GetConnection().
		Collection("notifications").Find(context.Background(), filter, opt)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err = cursor.Close(context.TODO()); err != nil {
			m.logger.Error(fmt.Sprintf("Resource release error: %s", err.Error()))
		}
	}()

	for cursor.Next(context.TODO()) {
		notif := new(model.Notification)
		if err = cursor.Decode(notif); err != nil {
			m.logger.Error(fmt.Sprintf("Decode error: %s", err.Error()))
			return nil, 0, err
		}

		result = append(result, *notif)
	}

	return result, total, nil
}

func NewMongoDbClient(db *infrastructure.Database, logger *zap.Logger) notification.Repository {
	return &MongoDbClient{
		db:     db,
		logger: logger,
	}
}
