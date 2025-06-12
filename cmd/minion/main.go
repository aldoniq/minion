package main

import (
	"fmt"
	"log"
	"os"

	"minion/internal/commands"
	"minion/internal/config"

	"github.com/spf13/cobra"
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

	// Создаем корневую команду
	rootCmd := &cobra.Command{
		Use:     "minion",
		Short:   "🍌 BELLO! Minion - Инструмент автоматизации для iiko API",
		Long:    "🍌 BELLO! Минион для автоматизации работы с iiko API. Продление ключей и обновление меню.",
		Version: Version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Показываем конфигурацию перед выполнением команд
			config.PrintEnvConfig(envConfig)
			fmt.Println()
		},
	}

	// Добавляем команды
	rootCmd.AddCommand(commands.ExtendKeysCmd)
	rootCmd.AddCommand(commands.RefreshMenusCmd)

	// Настройка версии
	rootCmd.SetVersionTemplate("🍌 Minion v{{.Version}} - BANANA!\n")

	// Выполняем команду
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
