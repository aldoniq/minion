package main

import (
	"fmt"
	"log"
	"os"

	"minion/internal/commands"
	"minion/internal/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "minion",
	Short:   "🍌 Minion - iiko Automation Tool",
	Long:    "🍌 BELLO! Minion - Инструмент автоматизации для iiko API",
	Version: "1.0.0",
}

func init() {
	rootCmd.AddCommand(commands.ExtendKeysCmd)
	rootCmd.AddCommand(commands.RefreshMenusCmd)
}

func main() {
	// Проверяем наличие конфигурационного файла
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		fmt.Println("⚠️  Файл config.json не найден. Создание образца...")
		if err := config.CreateSampleConfig(); err != nil {
			log.Fatal("❌ Ошибка создания образца конфигурации:", err)
		}
		fmt.Println("✅ Образец config.json создан. Пожалуйста, настройте параметры ресторанов.")
		return
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
