package database

import (
	"context"
	"fmt"
	"time"

	"minion/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RestaurantService предоставляет методы для работы с ресторанами в MongoDB
type RestaurantService struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

// NewRestaurantService создает новый экземпляр RestaurantService
func NewRestaurantService(connectionString, databaseName string) (*RestaurantService, error) {
	// Настройки подключения
	clientOptions := options.Client().ApplyURI(connectionString)

	// Создаем контекст с таймаутом для подключения
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Подключаемся к MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к MongoDB: %v", err)
	}

	// Проверяем соединение
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ошибка ping MongoDB: %v", err)
	}

	// Получаем базу данных и коллекцию
	database := client.Database(databaseName)
	collection := database.Collection("restaurants")

	return &RestaurantService{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

// GetActiveIikoRestaurants получает все активные рестораны с типом iiko
func (rs *RestaurantService) GetActiveIikoRestaurants() ([]*models.RestaurantMongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Фильтр для поиска активных iiko ресторанов
	filter := bson.M{
		"pos_type":   "iiko",
		"is_deleted": bson.M{"$ne": true},
		"$or": []bson.M{
			{"settings.is_deleted": bson.M{"$ne": true}},
			{"settings.is_deleted": bson.M{"$exists": false}},
		},
		// Проверяем что у ресторана есть данные iiko_cloud
		"iiko_cloud.iiko_web_domain": bson.M{"$ne": ""},
	}

	// Выполняем поиск
	cursor, err := rs.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска ресторанов: %v", err)
	}
	defer cursor.Close(ctx)

	// Декодируем результаты
	var restaurants []*models.RestaurantMongo
	if err := cursor.All(ctx, &restaurants); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ресторанов: %v", err)
	}

	return restaurants, nil
}

// Close закрывает соединение с базой данных
func (rs *RestaurantService) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return rs.client.Disconnect(ctx)
}
