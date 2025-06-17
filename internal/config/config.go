package config

import (
	"context"
	"os"

	"minion/internal/aws"
	"minion/internal/database"
	"minion/internal/models"
)

// EnvConfig содержит конфигурацию из переменных окружения
type EnvConfig struct {
	// Настройки HTTP сервера
	HTTPPort string // HTTP_PORT

	// Настройки AWS Secrets Manager
	AWSRegion     string // AWS_REGION
	AWSSecretName string // AWS_SECRET_NAME
}

// LoadEnvConfig загружает конфигурацию из переменных окружения
func LoadEnvConfig() *EnvConfig {
	return &EnvConfig{
		// Настройки HTTP сервера
		HTTPPort: getEnvWithDefault("HTTP_PORT", "3000"),

		// Настройки AWS Secrets Manager
		AWSRegion:     getEnvWithDefault("AWS_REGION", "eu-west-1"),
		AWSSecretName: getEnvWithDefault("AWS_SECRET_NAME", "ProdEnvs"),
	}
}

// LoadRestaurants загружает рестораны из базы данных через AWS Secrets Manager
func LoadRestaurants(ctx context.Context, envConfig *EnvConfig) ([]*models.Restaurant, error) {
	// Создаем AWS Secrets Manager клиент
	awsClient, err := aws.NewSecretsManager(envConfig.AWSRegion)
	if err != nil {
		return nil, err
	}

	// Получаем credentials из AWS
	dbCredentials, err := awsClient.GetDatabaseCredentials(envConfig.AWSSecretName)
	if err != nil {
		return nil, err
	}

	// Создаем сервис для работы с ресторанами
	restaurantService, err := database.NewRestaurantService(dbCredentials.DbURL, dbCredentials.DbName)
	if err != nil {
		return nil, err
	}
	defer restaurantService.Close()

	// Загружаем рестораны из базы данных
	mongoRestaurants, err := restaurantService.GetActiveIikoRestaurants()
	if err != nil {
		return nil, err
	}

	// Конвертируем в формат для minion
	var restaurants []*models.Restaurant
	for _, mongoRestaurant := range mongoRestaurants {
		restaurant := mongoRestaurant.ToMinion()
		if restaurant != nil {
			restaurants = append(restaurants, restaurant)
		}
	}

	return restaurants, nil
}

// getEnvWithDefault получает значение переменной окружения или возвращает значение по умолчанию
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
