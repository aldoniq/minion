package config

import (
	"encoding/json"
	"os"

	"minion/internal/models"
)

// LoadConfig загружает конфигурацию из файла
func LoadConfig(filename string) (*models.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Устанавливаем значение по умолчанию для продления ключей
	if config.ExtensionYears == 0 {
		config.ExtensionYears = 2
	}

	return &config, nil
}

// CreateSampleConfig создает образец конфигурационного файла
func CreateSampleConfig() error {
	sampleConfig := models.Config{
		ExtensionYears: 2,
		Restaurants: []models.Restaurant{
			{
				Name:     "Казбек Бокейхана",
				BaseURL:  "https://kazbek-bokeihana.iikoweb.ru",
				Login:    "2020",
				Password: "2020",
				Enabled:  true,
			},
		},
	}

	data, err := json.MarshalIndent(sampleConfig, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("config.json", data, 0644)
}
