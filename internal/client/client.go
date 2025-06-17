package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"minion/internal/models"
)

// IikoClient - HTTP клиент для работы с iiko API
type IikoClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewIikoClient создает новый экземпляр клиента
func NewIikoClient(baseURL string) *IikoClient {
	return &IikoClient{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Login выполняет авторизацию и возвращает PHPSESSID
func (c *IikoClient) Login(login, password string) (string, error) {
	loginData := models.LoginRequest{Login: login, Password: password}
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

// GetApiLogins получает список API логинов
func (c *IikoClient) GetApiLogins(sessionID string) (*models.ApiLoginsResponse, error) {
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

	var response models.ApiLoginsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetApiLoginDetail получает детальную информацию об API логине
func (c *IikoClient) GetApiLoginDetail(sessionID string, apiLoginID string) (*models.ApiLoginDetailResponse, error) {
	requestData := models.ApiLoginDetailRequest{ApiLoginID: apiLoginID}
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

	var response models.ApiLoginDetailResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SaveApiLoginDetail сохраняет обновленный API логин
func (c *IikoClient) SaveApiLoginDetail(sessionID string, apiLoginDetail models.ApiLoginDetail) error {
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

// GetExternalMenus получает список внешних меню
func (c *IikoClient) GetExternalMenus(sessionID string) (*models.ExternalMenuResponse, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/api/external-menu", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "ru_RU")
	req.Header.Set("Cookie", "PHPSESSID="+sessionID)
	req.Header.Set("Referer", c.baseURL+"/external-menu/index.html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка получения меню, статус: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response models.ExternalMenuResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// RefreshExternalMenu обновляет внешнее меню
func (c *IikoClient) RefreshExternalMenu(sessionID string, menuID int) error {
	refreshData := models.RefreshMenuRequest{
		RefreshNameAndDescription:       false,
		RefreshPrice:                    true,
		RefreshImages:                   false,
		RefreshModifiersNumber:          true,
		RefreshNutritionPerHundredGrams: true,
		RefreshAllergens:                true,
		RefreshCombos:                   true,
	}

	jsonData, _ := json.Marshal(refreshData)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/external-menu/refresh-menu/%d", c.baseURL, menuID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "ru_RU")
	req.Header.Set("Cookie", "PHPSESSID="+sessionID)
	req.Header.Set("Referer", c.baseURL+"/external-menu/index.html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка обновления меню, статус: %d", resp.StatusCode)
	}

	return nil
}
