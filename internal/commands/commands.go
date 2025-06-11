package commands

import (
	"fmt"
	"log"
	"strings"
	"time"

	"minion/internal/client"
	"minion/internal/config"
	"minion/internal/models"

	"github.com/spf13/cobra"
)

// ExtendKeysCmd - команда продления API ключей
var ExtendKeysCmd = &cobra.Command{
	Use:   "extend-keys",
	Short: "Продление срока действия API ключей",
	Long:  "🔑 Продление срока действия API ключей до максимального значения",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🍌 BELLO! Запуск продления API ключей...")

		cfg, err := config.LoadConfig("config.json")
		if err != nil {
			log.Fatal("❌ Ошибка загрузки конфигурации:", err)
		}

		totalUpdated := 0
		successCount := 0
		failedCount := 0

		for _, restaurant := range cfg.Restaurants {
			if !restaurant.Enabled {
				fmt.Printf("⏭️  Ресторан %s отключен, пропускаем\n", restaurant.Name)
				continue
			}

			updated, err := processExtendKeys(restaurant, cfg.ExtensionYears)
			if err != nil {
				fmt.Printf("❌ Ошибка обработки ресторана %s: %v\n", restaurant.Name, err)
				failedCount++
			} else {
				fmt.Printf("🎉 GELATO! Ресторан %s: обновлено %d ключей\n", restaurant.Name, updated)
				totalUpdated += updated
				successCount++
			}
		}

		printSummary("API ключей", successCount, failedCount, totalUpdated)
	},
}

// RefreshMenusCmd - команда обновления меню
var RefreshMenusCmd = &cobra.Command{
	Use:   "refresh-menus",
	Short: "Обновление внешних меню",
	Long:  "🍽️ Обновление внешних меню для всех ресторанов",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🍌 BELLO! Запуск обновления меню...")

		cfg, err := config.LoadConfig("config.json")
		if err != nil {
			log.Fatal("❌ Ошибка загрузки конфигурации:", err)
		}

		totalUpdated := 0
		successCount := 0
		failedCount := 0

		for _, restaurant := range cfg.Restaurants {
			if !restaurant.Enabled {
				fmt.Printf("⏭️  Ресторан %s отключен, пропускаем\n", restaurant.Name)
				continue
			}

			updated, err := processRefreshMenus(restaurant)
			if err != nil {
				fmt.Printf("❌ Ошибка обработки ресторана %s: %v\n", restaurant.Name, err)
				failedCount++
			} else {
				fmt.Printf("🎉 GELATO! Ресторан %s: обновлено %d меню\n", restaurant.Name, updated)
				totalUpdated += updated
				successCount++
			}
		}

		printSummary("меню", successCount, failedCount, totalUpdated)
	},
}

// processExtendKeys обрабатывает продление ключей для одного ресторана
func processExtendKeys(restaurant models.Restaurant, extensionYears int) (int, error) {
	fmt.Printf("\n🔄 Обработка ресторана: %s\n", restaurant.Name)

	apiClient := client.NewIikoClient(restaurant.BaseURL)

	// Авторизация
	sessionID, err := apiClient.Login(restaurant.Login, restaurant.Password)
	if err != nil {
		return 0, fmt.Errorf("ошибка авторизации: %v", err)
	}
	fmt.Println("✅ Авторизация успешна!")

	// Получение API логинов
	response, err := apiClient.GetApiLogins(sessionID)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения API логинов: %v", err)
	}
	fmt.Printf("📋 Найдено %d API логинов\n", len(response.ApiLogins))

	updatedCount := 0
	for _, apiLogin := range response.ApiLogins {
		if !apiLogin.IsActive {
			fmt.Printf("⏭️  API логин %s неактивен, пропускаем\n", apiLogin.Name)
			continue
		}

		fmt.Printf("🔍 Обработка API логина: %s\n", apiLogin.Name)

		// Получаем детальную информацию
		detailResponse, err := apiClient.GetApiLoginDetail(sessionID, apiLogin.ID)
		if err != nil {
			fmt.Printf("⚠️  Ошибка получения деталей: %v\n", err)
			continue
		}

		newExpirationDate, err := extendExpirationDate(detailResponse.ApiLoginInfo.ExpirationDate, extensionYears)
		if err != nil {
			fmt.Printf("⚠️  Ошибка обработки даты: %v\n", err)
			continue
		}

		fmt.Printf("📅 Текущий срок: %s\n", detailResponse.ApiLoginInfo.ExpirationDate)

		if newExpirationDate == detailResponse.ApiLoginInfo.ExpirationDate {
			fmt.Printf("✅ Срок действия уже максимальный, пропускаем\n")
			continue
		}

		fmt.Printf("📅 Новый срок: %s\n", newExpirationDate)

		// Обновляем дату в детальной структуре
		detailResponse.ApiLoginInfo.ExpirationDate = newExpirationDate

		err = apiClient.SaveApiLoginDetail(sessionID, detailResponse.ApiLoginInfo)
		if err != nil {
			fmt.Printf("⚠️  Ошибка сохранения: %v\n", err)
			continue
		}

		fmt.Printf("✅ API логин %s успешно обновлен!\n", apiLogin.Name)
		updatedCount++
	}

	return updatedCount, nil
}

// processRefreshMenus обрабатывает обновление меню для одного ресторана
func processRefreshMenus(restaurant models.Restaurant) (int, error) {
	fmt.Printf("\n🔄 Обработка ресторана: %s\n", restaurant.Name)

	apiClient := client.NewIikoClient(restaurant.BaseURL)

	// Авторизация
	sessionID, err := apiClient.Login(restaurant.Login, restaurant.Password)
	if err != nil {
		return 0, fmt.Errorf("ошибка авторизации: %v", err)
	}
	fmt.Println("✅ Авторизация успешна!")

	// Получение списка внешних меню
	menuResponse, err := apiClient.GetExternalMenus(sessionID)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения списка меню: %v", err)
	}

	if menuResponse.Error {
		return 0, fmt.Errorf("ошибка в ответе API")
	}

	fmt.Printf("📋 Найдено %d внешних меню\n", len(menuResponse.Data))

	updatedCount := 0
	for _, menu := range menuResponse.Data {
		fmt.Printf("🔍 Обновление меню: %s (ID: %d)\n", menu.Name, menu.ID)

		err = apiClient.RefreshExternalMenu(sessionID, menu.ID)
		if err != nil {
			fmt.Printf("⚠️  Ошибка обновления меню %s: %v\n", menu.Name, err)
			continue
		}

		fmt.Printf("✅ Меню %s успешно обновлено!\n", menu.Name)
		updatedCount++
	}

	return updatedCount, nil
}

// extendExpirationDate продлевает срок действия даты
func extendExpirationDate(currentDate string, years int) (string, error) {
	formats := []string{
		"2006-01-02 15:04:05.000",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	}

	var parsedTime time.Time
	var err error

	for _, format := range formats {
		parsedTime, err = time.Parse(format, currentDate)
		if err == nil {
			break
		}
	}

	if err != nil {
		return "", fmt.Errorf("не удалось распарсить дату %s: %v", currentDate, err)
	}

	// Максимальная дата = сегодня + указанное количество лет
	maxDate := time.Now().AddDate(years, 0, 0)

	// Если текущая дата уже максимальная или больше, не изменяем
	if !parsedTime.Before(maxDate) {
		return currentDate, nil
	}

	// Устанавливаем максимальную дату
	return maxDate.Format("2006-01-02 15:04:05.000"), nil
}

// printSummary выводит итоговый отчет
func printSummary(itemType string, successCount, failedCount, totalUpdated int) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📊 ИТОГОВЫЙ ОТЧЕТ")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("✅ Успешно обработано ресторанов: %d\n", successCount)
	fmt.Printf("❌ Ошибок: %d\n", failedCount)
	fmt.Printf("🔄 Всего обновлено %s: %d\n", itemType, totalUpdated)

	if failedCount == 0 {
		fmt.Println("\n🍌 BANANA! Все рестораны успешно обработаны! 🎉")
	} else {
		fmt.Println("\n⚠️  Обработка завершена с ошибками")
	}
}
