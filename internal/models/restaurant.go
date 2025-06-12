package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RestaurantMongo представляет документ ресторана в MongoDB
type RestaurantMongo struct {
	ID                          primitive.ObjectID `bson:"_id,omitempty"`
	Token                       string             `bson:"token"`
	Name                        string             `bson:"name"`
	PosType                     string             `bson:"pos_type"`
	City                        string             `bson:"city"`
	StorePhoneNumber            string             `bson:"store_phone_number"`
	SendWhatsappNotification    bool               `bson:"send_whatsapp_notification"`
	WhatsappErrorStoplistChatID string             `bson:"whatsapp_error_stoplist_chat_id"`
	IikoCloud                   IikoCloudConfig    `bson:"iiko_cloud"`
	Settings                    RestaurantSettings `bson:"settings"`
	SendToPos                   bool               `bson:"send_to_pos"`
	IsDeleted                   bool               `bson:"is_deleted"`
	IntegrationDate             time.Time          `bson:"integration_date"`
	UpdatedAt                   time.Time          `bson:"updated_at"`
	CreatedAt                   time.Time          `bson:"created_at"`
}

// IikoCloudConfig содержит настройки для iiko Cloud
type IikoCloudConfig struct {
	OrganizationID string `bson:"organization_id"`
	TerminalID     string `bson:"terminal_id"`
	Key            string `bson:"key"`
	Login          string `bson:"login"`
	Password       string `bson:"password"`
	IsExternalMenu bool   `bson:"is_external_menu"`
	ExternalMenuID string `bson:"external_menu_id"`
	CustomDomain   string `bson:"custom_domain"`
}

// RestaurantSettings содержит общие настройки ресторана
type RestaurantSettings struct {
	SendToPos     bool   `bson:"send_to_pos"`
	IsMarketplace bool   `bson:"is_marketplace"`
	IsDeleted     bool   `bson:"is_deleted"`
	LanguageCode  string `bson:"language_code"`
}

// ToMinion конвертирует RestaurantMongo в Restaurant для minion
func (r *RestaurantMongo) ToMinion() *Restaurant {
	// Проверяем, что это iiko ресторан и у него есть необходимые данные
	if r.PosType != "iiko" || r.IikoCloud.CustomDomain == "" {
		return nil
	}

	// Формируем базовый URL
	baseURL := "https://" + r.IikoCloud.CustomDomain

	// Создаем Restaurant для minion
	return &Restaurant{
		Name:     r.Name,
		BaseURL:  baseURL,
		Login:    r.IikoCloud.Login,
		Password: r.IikoCloud.Password,
		Enabled:  !r.IsDeleted && !r.Settings.IsDeleted,
	}
}

// DatabaseCredentials представляет данные для подключения к базе
type DatabaseCredentials struct {
	DbURL    string `json:"db_url"`
	DbName   string `json:"db_name"`
	DbEngine string `json:"db_engine"`
}
