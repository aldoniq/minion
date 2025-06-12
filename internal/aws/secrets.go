package aws

import (
	"encoding/json"
	"fmt"

	"minion/internal/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// SecretsManager предоставляет интерфейс для работы с AWS Secrets Manager
type SecretsManager struct {
	client *secretsmanager.SecretsManager
	region string
}

// NewSecretsManager создает новый экземпляр SecretsManager
func NewSecretsManager(region string) (*SecretsManager, error) {
	// Создаем AWS сессию
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка создания AWS сессии: %v", err)
	}

	// Создаем клиент Secrets Manager
	client := secretsmanager.New(sess)

	return &SecretsManager{
		client: client,
		region: region,
	}, nil
}

// GetDatabaseCredentials получает данные для подключения к базе данных
func (sm *SecretsManager) GetDatabaseCredentials(secretName string) (*models.DatabaseCredentials, error) {
	// Получаем секрет
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := sm.client.GetSecretValue(input)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения секрета %s: %v", secretName, err)
	}

	// Парсим JSON
	var credentials models.DatabaseCredentials
	if err := json.Unmarshal([]byte(*result.SecretString), &credentials); err != nil {
		return nil, fmt.Errorf("ошибка парсинга секрета %s: %v", secretName, err)
	}

	return &credentials, nil
}

// GetSecretValue получает произвольное значение секрета
func (sm *SecretsManager) GetSecretValue(secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := sm.client.GetSecretValue(input)
	if err != nil {
		return "", fmt.Errorf("ошибка получения секрета %s: %v", secretName, err)
	}

	if result.SecretString == nil {
		return "", fmt.Errorf("секрет %s пуст", secretName)
	}

	return *result.SecretString, nil
}
