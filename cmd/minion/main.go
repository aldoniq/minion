package main

import (
	"fmt"
	"log"
	"os"

	"minion/internal/config"
	"minion/internal/server"
)

var Version = "2.1.0"

func main() {
	// Загружаем .env файл, если не найден - падаем
	if err := config.LoadEnvFile(); err != nil {
		fmt.Printf("❌ %v\n", err)
		fmt.Println("💡 Создайте .env файл на основе env.example")
		os.Exit(1)
	}

	// Загружаем и валидируем конфигурацию из переменных окружения
	envConfig := config.LoadEnvConfig()

	// Валидируем конфигурацию
	if errors := config.ValidateEnvConfig(envConfig); len(errors) > 0 {
		fmt.Println("❌ Ошибки конфигурации:")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
		fmt.Println("\n💡 Проверьте переменные в .env файле")
		os.Exit(1)
	}

	// Показываем конфигурацию
	config.PrintEnvConfig(envConfig)
	fmt.Println()

	// Запускаем HTTP сервер
	fmt.Println("🍌 BELLO! Запуск Minion HTTP API сервера...")
	if err := server.StartServer(envConfig.HTTPPort); err != nil {
		log.Fatalf("❌ Ошибка запуска сервера: %v", err)
	}
}
