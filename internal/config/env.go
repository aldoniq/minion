package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

// LoadEnvFile загружает .env файл, если не найден - ошибка
func LoadEnvFile() error {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("файл .env не найден или содержит ошибки: %v", err)
	}
	return nil
}

// ValidateEnvConfig проверяет корректность конфигурации из переменных окружения
func ValidateEnvConfig(config *EnvConfig) []string {
	var errors []string

	// AWS Secrets Manager
	if config.AWSRegion == "" {
		errors = append(errors, "AWS_REGION не может быть пустой")
	}
	if config.AWSSecretName == "" {
		errors = append(errors, "AWS_SECRET_NAME не может быть пустой")
	}

	return errors
}

// PrintEnvConfig выводит текущую конфигурацию (без секретных данных)
func PrintEnvConfig(config *EnvConfig) {
	fmt.Println("🔧 Текущая конфигурация:")
	fmt.Printf("  🚀 HTTP Port: %s\n", config.HTTPPort)
	fmt.Printf("  🌍 AWS Region: %s\n", config.AWSRegion)
	fmt.Printf("  🔑 AWS Secret Name: %s\n", config.AWSSecretName)
}
