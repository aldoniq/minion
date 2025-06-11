package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Структуры данных
type Config struct {
	ExtensionYears int          `json:"extension_years"`
	Restaurants    []Restaurant `json:"restaurants"`
}

type Restaurant struct {
	Name     string `json:"name"`
	BaseURL  string `json:"base_url"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Enabled  bool   `json:"enabled"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ApiLogin struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	IsActive       bool   `json:"isActive"`
	ExpirationDate string `json:"expirationDate"`
}

type ApiLoginsResponse struct {
	CorrelationID string     `json:"correlationId"`
	ApiLogins     []ApiLogin `json:"apiLogins"`
}

// Полная структура для детального API логина
type ApiLoginDetail struct {
	ID              string         `json:"id"`
	APIKey          string         `json:"apiKey"`
	Name            string         `json:"name"`
	SourceKey       string         `json:"sourceKey"`
	IsActive        bool           `json:"isActive"`
	IncludedRmses   []IncludedRms  `json:"includedRmses"`
	TemplateID      string         `json:"templateId"`
	ExternalMenus   []ExternalMenu `json:"externalMenus"`
	PriceCategories []interface{}  `json:"priceCategories"`
	Email           string         `json:"email"`
	ExpirationDate  string         `json:"expirationDate"`
	IsLongLived     bool           `json:"isLongLived"`
}

type IncludedRms struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	Rating            string          `json:"rating"`
	CrmID             interface{}     `json:"crmId"`
	WebHookSettings   WebHookSettings `json:"webHookSettings"`
	LastConnectToUoc  string          `json:"lastConnectToUoc"`
	ConnectedToCard   bool            `json:"connectedToCard"`
	IsCloud           bool            `json:"isCloud"`
	HasRestApiLicense bool            `json:"hasRestApiLicense"`
}

type WebHookSettings struct {
	Blocked        bool           `json:"blocked"`
	WebHooksUri    string         `json:"webHooksUri"`
	AuthToken      string         `json:"authToken"`
	WebHooksFilter WebHooksFilter `json:"webHooksFilter"`
}

type WebHooksFilter struct {
	DeliveryOrderFilter      OrderFilter   `json:"deliveryOrderFilter"`
	TableOrderFilter         OrderFilter   `json:"tableOrderFilter"`
	ReserveFilter            ReserveFilter `json:"reserveFilter"`
	StopListUpdateFilter     UpdateFilter  `json:"stopListUpdateFilter"`
	PersonalShiftFilter      UpdateFilter  `json:"personalShiftFilter"`
	NomenclatureUpdateFilter UpdateFilter  `json:"nomenclatureUpdateFilter"`
}

type OrderFilter struct {
	OrderStatuses []string `json:"orderStatuses"`
	ItemStatuses  []string `json:"itemStatuses"`
	Errors        bool     `json:"errors"`
}

type ReserveFilter struct {
	Updates bool `json:"updates"`
	Errors  bool `json:"errors"`
}

type UpdateFilter struct {
	Updates bool `json:"updates"`
}

type ExternalMenu struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ApiLoginDetailRequest struct {
	ApiLoginID string `json:"apiLoginId"`
}

type ApiLoginDetailResponse struct {
	CorrelationID string         `json:"correlationId"`
	ApiLoginInfo  ApiLoginDetail `json:"apiLoginInfo"`
}

// HTTP клиент для iiko API
type IikoClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewIikoClient(baseURL string) *IikoClient {
	return &IikoClient{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Авторизация и получение PHPSESSID
func (c *IikoClient) Login(login, password string) (string, error) {
	loginData := LoginRequest{Login: login, Password: password}
	jsonData, _ := json.Marshal(loginData)

	req, err := http.NewRequest("POST", c.baseURL+"/api/auth/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("авторизация не удалась, статус: %d", resp.StatusCode)
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "PHPSESSID" {
			return cookie.Value, nil
		}
	}

	return "", fmt.Errorf("PHPSESSID не найден")
}

// Получение списка API логинов
func (c *IikoClient) GetApiLogins(sessionID string) (*ApiLoginsResponse, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/api/integration-management/api-logins/get-all", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "ru_RU")
	req.Header.Set("Cookie", "PHPSESSID="+sessionID)
	req.Header.Set("Origin", c.baseURL)
	req.Header.Set("Referer", c.baseURL+"/integration-management/index.html")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка получения логинов, статус: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response ApiLoginsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Продление срока действия даты
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

// Загрузка конфигурации
func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.ExtensionYears == 0 {
		config.ExtensionYears = 2
	}

	return &config, nil
}

// Обработка одного ресторана
func processRestaurant(restaurant Restaurant, extensionYears int) (int, error) {
	fmt.Printf("\n🔄 Обработка ресторана: %s\n", restaurant.Name)

	client := NewIikoClient(restaurant.BaseURL)

	// Авторизация
	sessionID, err := client.Login(restaurant.Login, restaurant.Password)
	if err != nil {
		return 0, fmt.Errorf("ошибка авторизации: %v", err)
	}
	fmt.Println("✅ Авторизация успешна!")

	// Получение API логинов
	response, err := client.GetApiLogins(sessionID)
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
		detailResponse, err := client.GetApiLoginDetail(sessionID, apiLogin.ID)
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

		err = client.SaveApiLoginDetail(sessionID, detailResponse.ApiLoginInfo)
		if err != nil {
			fmt.Printf("⚠️  Ошибка сохранения: %v\n", err)
			continue
		}

		fmt.Printf("✅ API логин %s успешно обновлен!\n", apiLogin.Name)
		updatedCount++
	}

	return updatedCount, nil
}

// Получение детальной информации об API логине
func (c *IikoClient) GetApiLoginDetail(sessionID string, apiLoginID string) (*ApiLoginDetailResponse, error) {
	requestData := ApiLoginDetailRequest{ApiLoginID: apiLoginID}
	jsonData, _ := json.Marshal(requestData)

	req, err := http.NewRequest("POST", c.baseURL+"/api/integration-management/api-logins/get", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "ru_RU")
	req.Header.Set("Cookie", "PHPSESSID="+sessionID)
	req.Header.Set("Origin", c.baseURL)
	req.Header.Set("Referer", c.baseURL+"/integration-management/index.html")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка получения деталей, статус: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response ApiLoginDetailResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Сохранение обновленного API логина
func (c *IikoClient) SaveApiLoginDetail(sessionID string, apiLoginDetail ApiLoginDetail) error {
	jsonData, _ := json.Marshal(apiLoginDetail)

	req, err := http.NewRequest("POST", c.baseURL+"/api/integration-management/save-api-login", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "ru_RU")
	req.Header.Set("Cookie", "PHPSESSID="+sessionID)
	req.Header.Set("Origin", c.baseURL)
	req.Header.Set("Referer", c.baseURL+"/integration-management/index.html")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка сохранения, статус: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	fmt.Println("🍌 BELLO! Minion API Key Extension Tool запущен!")

	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal("❌ Ошибка загрузки конфигурации:", err)
	}

	totalUpdated := 0
	successCount := 0
	failedCount := 0

	for _, restaurant := range config.Restaurants {
		if !restaurant.Enabled {
			fmt.Printf("⏭️  Ресторан %s отключен, пропускаем\n", restaurant.Name)
			continue
		}

		updated, err := processRestaurant(restaurant, config.ExtensionYears)
		if err != nil {
			fmt.Printf("❌ Ошибка обработки ресторана %s: %v\n", restaurant.Name, err)
			failedCount++
		} else {
			fmt.Printf("🎉 GELATO! Ресторан %s: обновлено %d ключей\n", restaurant.Name, updated)
			totalUpdated += updated
			successCount++
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📊 ИТОГОВЫЙ ОТЧЕТ")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("✅ Успешно обработано ресторанов: %d\n", successCount)
	fmt.Printf("❌ Ошибок: %d\n", failedCount)
	fmt.Printf("🔑 Всего продлено API ключей: %d\n", totalUpdated)

	if failedCount == 0 {
		fmt.Println("\n🍌 BANANA! Все рестораны успешно обработаны! 🎉")
	} else {
		fmt.Println("\n⚠️  Обработка завершена с ошибками")
	}
}
