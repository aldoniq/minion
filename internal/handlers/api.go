package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"minion/internal/client"
	"minion/internal/config"
	"minion/internal/models"

	"github.com/gofiber/fiber/v2"
)

// APIResponse представляет стандартный ответ API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// OperationResult содержит результаты выполнения операции
type OperationResult struct {
	ProcessedRestaurants int                `json:"processed_restaurants"`
	Successful           int                `json:"successful"`
	Failed               int                `json:"failed"`
	Duration             string             `json:"duration"`
	Details              []RestaurantResult `json:"details"`
}

// RestaurantResult содержит результат обработки одного ресторана
type RestaurantResult struct {
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Updated int    `json:"updated"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// HealthCheck проверка состояния сервиса
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(APIResponse{
		Success: true,
		Message: "🍌 BELLO! Minion API работает",
		Data: fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "2.1.0",
		},
	})
}

// GetConfig показывает текущую конфигурацию
func GetConfig(c *fiber.Ctx) error {
	envConfig := config.LoadEnvConfig()

	return c.JSON(APIResponse{
		Success: true,
		Message: "🔧 Текущая конфигурация",
		Data: fiber.Map{
			"aws_region":      envConfig.AWSRegion,
			"aws_secret_name": envConfig.AWSSecretName,
		},
	})
}

// ExtendKeys обработчик продления API ключей
func ExtendKeys(c *fiber.Ctx) error {
	log.Printf("🔑 API запрос: продление ключей от %s", c.IP())

	startTime := time.Now()

	// Загружаем рестораны
	restaurants, extensionYears, err := loadRestaurants()
	if err != nil {
		log.Printf("❌ Ошибка загрузки ресторанов: %v", err)
		return c.Status(500).JSON(APIResponse{
			Success: false,
			Message: "Ошибка загрузки ресторанов",
			Error:   err.Error(),
		})
	}

	result := OperationResult{
		ProcessedRestaurants: len(restaurants),
		Details:              make([]RestaurantResult, 0),
	}

	// Обрабатываем каждый ресторан
	for _, restaurant := range restaurants {
		if !restaurant.Enabled {
			log.Printf("⏭️  Ресторан %s отключен, пропускаем", restaurant.Name)
			continue
		}

		restaurantResult := RestaurantResult{
			Name: restaurant.Name,
		}

		updated, err := processExtendKeys(*restaurant, extensionYears)
		if err != nil {
			log.Printf("❌ Ошибка обработки ресторана %s: %v", restaurant.Name, err)
			restaurantResult.Success = false
			restaurantResult.Error = err.Error()
			result.Failed++
		} else {
			log.Printf("✅ Ресторан %s: обновлено %d ключей", restaurant.Name, updated)
			restaurantResult.Success = true
			restaurantResult.Updated = updated
			restaurantResult.Message = fmt.Sprintf("Обновлено %d ключей", updated)
			result.Successful++
		}

		result.Details = append(result.Details, restaurantResult)
	}

	result.Duration = time.Since(startTime).String()

	log.Printf("🎉 GELATO! Продление ключей завершено: %d успешно, %d ошибок",
		result.Successful, result.Failed)

	return c.JSON(APIResponse{
		Success: true,
		Message: fmt.Sprintf("🎉 GELATO! Обработано %d ресторанов", result.ProcessedRestaurants),
		Data:    result,
	})
}

// RefreshMenus обработчик обновления меню
func RefreshMenus(c *fiber.Ctx) error {
	log.Printf("🍽️ API запрос: обновление меню от %s", c.IP())

	startTime := time.Now()

	// Загружаем рестораны
	restaurants, _, err := loadRestaurants()
	if err != nil {
		log.Printf("❌ Ошибка загрузки ресторанов: %v", err)
		return c.Status(500).JSON(APIResponse{
			Success: false,
			Message: "Ошибка загрузки ресторанов",
			Error:   err.Error(),
		})
	}

	result := OperationResult{
		ProcessedRestaurants: len(restaurants),
		Details:              make([]RestaurantResult, 0),
	}

	// Обрабатываем каждый ресторан
	for _, restaurant := range restaurants {
		if !restaurant.Enabled {
			log.Printf("⏭️  Ресторан %s отключен, пропускаем", restaurant.Name)
			continue
		}

		restaurantResult := RestaurantResult{
			Name: restaurant.Name,
		}

		updated, err := processRefreshMenus(*restaurant)
		if err != nil {
			log.Printf("❌ Ошибка обработки ресторана %s: %v", restaurant.Name, err)
			restaurantResult.Success = false
			restaurantResult.Error = err.Error()
			result.Failed++
		} else {
			log.Printf("✅ Ресторан %s: обновлено %d меню", restaurant.Name, updated)
			restaurantResult.Success = true
			restaurantResult.Updated = updated
			restaurantResult.Message = fmt.Sprintf("Обновлено %d меню", updated)
			result.Successful++
		}

		result.Details = append(result.Details, restaurantResult)
	}

	result.Duration = time.Since(startTime).String()

	log.Printf("🎉 GELATO! Обновление меню завершено: %d успешно, %d ошибок",
		result.Successful, result.Failed)

	return c.JSON(APIResponse{
		Success: true,
		Message: fmt.Sprintf("🎉 GELATO! Обработано %d ресторанов", result.ProcessedRestaurants),
		Data:    result,
	})
}

// loadRestaurants загружает рестораны из базы данных
func loadRestaurants() ([]*models.Restaurant, int, error) {
	envConfig := config.LoadEnvConfig()
	restaurants, err := config.LoadRestaurants(context.Background(), envConfig)
	if err != nil {
		return nil, 0, err
	}

	// Возвращаем дефолтное значение для extension_years
	return restaurants, 2, nil
}

// processExtendKeys обрабатывает продление ключей для одного ресторана
func processExtendKeys(restaurant models.Restaurant, extensionYears int) (int, error) {
	apiClient := client.NewIikoClient(restaurant.BaseURL)

	// Авторизация
	sessionID, err := apiClient.Login(restaurant.Login, restaurant.Password)
	if err != nil {
		return 0, fmt.Errorf("ошибка авторизации: %v", err)
	}

	// Получение API логинов
	response, err := apiClient.GetApiLogins(sessionID)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения API логинов: %v", err)
	}

	updatedCount := 0
	for _, apiLogin := range response.ApiLogins {
		if !apiLogin.IsActive {
			continue
		}

		// Получаем детальную информацию
		detailResponse, err := apiClient.GetApiLoginDetail(sessionID, apiLogin.ID)
		if err != nil {
			continue
		}

		newExpirationDate, err := extendExpirationDate(detailResponse.ApiLoginInfo.ExpirationDate, extensionYears)
		if err != nil {
			continue
		}

		if newExpirationDate == detailResponse.ApiLoginInfo.ExpirationDate {
			continue
		}

		// Обновляем дату
		detailResponse.ApiLoginInfo.ExpirationDate = newExpirationDate
		err = apiClient.SaveApiLoginDetail(sessionID, detailResponse.ApiLoginInfo)
		if err != nil {
			continue
		}

		updatedCount++
	}

	return updatedCount, nil
}

// processRefreshMenus обрабатывает обновление меню для одного ресторана
func processRefreshMenus(restaurant models.Restaurant) (int, error) {
	apiClient := client.NewIikoClient(restaurant.BaseURL)

	// Авторизация
	sessionID, err := apiClient.Login(restaurant.Login, restaurant.Password)
	if err != nil {
		return 0, fmt.Errorf("ошибка авторизации: %v", err)
	}

	// Получение списка внешних меню
	menus, err := apiClient.GetExternalMenus(sessionID)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения меню: %v", err)
	}

	updatedCount := 0
	for _, menu := range menus.Data {
		// Обновляем меню
		err := apiClient.RefreshExternalMenu(sessionID, menu.ID)
		if err != nil {
			continue
		}
		updatedCount++
	}

	return updatedCount, nil
}

// extendExpirationDate продлевает дату истечения на указанное количество лет
func extendExpirationDate(currentDate string, years int) (string, error) {
	// Парсим текущую дату
	parsedTime, err := time.Parse("02.01.2006", currentDate)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга даты %s: %v", currentDate, err)
	}

	// Добавляем годы
	newTime := parsedTime.AddDate(years, 0, 0)

	// Максимальная дата - 31.12.2099
	maxDate := time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC)
	if newTime.After(maxDate) {
		newTime = maxDate
	}

	return newTime.Format("02.01.2006"), nil
}
