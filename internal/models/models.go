package models

// Ресторан
type Restaurant struct {
	Name     string `json:"name"`
	BaseURL  string `json:"base_url"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Enabled  bool   `json:"enabled"`
}

// Запрос авторизации
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// API логин
type ApiLogin struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	IsActive       bool   `json:"isActive"`
	ExpirationDate string `json:"expirationDate"`
}

// Ответ со списком API логинов
type ApiLoginsResponse struct {
	CorrelationID string     `json:"correlationId"`
	ApiLogins     []ApiLogin `json:"apiLogins"`
}

// Детальная информация об API логине
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

// Включенные RMS
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

// Настройки WebHook
type WebHookSettings struct {
	Blocked        bool           `json:"blocked"`
	WebHooksUri    string         `json:"webHooksUri"`
	AuthToken      string         `json:"authToken"`
	WebHooksFilter WebHooksFilter `json:"webHooksFilter"`
}

// Фильтр WebHook
type WebHooksFilter struct {
	DeliveryOrderFilter      OrderFilter   `json:"deliveryOrderFilter"`
	TableOrderFilter         OrderFilter   `json:"tableOrderFilter"`
	ReserveFilter            ReserveFilter `json:"reserveFilter"`
	StopListUpdateFilter     UpdateFilter  `json:"stopListUpdateFilter"`
	PersonalShiftFilter      UpdateFilter  `json:"personalShiftFilter"`
	NomenclatureUpdateFilter UpdateFilter  `json:"nomenclatureUpdateFilter"`
}

// Фильтр заказов
type OrderFilter struct {
	OrderStatuses []string `json:"orderStatuses"`
	ItemStatuses  []string `json:"itemStatuses"`
	Errors        bool     `json:"errors"`
}

// Фильтр резервов
type ReserveFilter struct {
	Updates bool `json:"updates"`
	Errors  bool `json:"errors"`
}

// Фильтр обновлений
type UpdateFilter struct {
	Updates bool `json:"updates"`
}

// Внешнее меню
type ExternalMenu struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Запрос деталей API логина
type ApiLoginDetailRequest struct {
	ApiLoginID string `json:"apiLoginId"`
}

// Ответ с деталями API логина
type ApiLoginDetailResponse struct {
	CorrelationID string         `json:"correlationId"`
	ApiLoginInfo  ApiLoginDetail `json:"apiLoginInfo"`
}

// Ответ со списком внешних меню
type ExternalMenuResponse struct {
	Error               bool                 `json:"error"`
	Warning             bool                 `json:"warning"`
	Data                []ExternalMenuDetail `json:"data"`
	FormValidationError bool                 `json:"formValidationError"`
}

// Детали внешнего меню
type ExternalMenuDetail struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	GeneratingStatus int    `json:"generatingStatus"`
	PriceCategoryID  string `json:"priceCategoryId"`
}

// Запрос обновления меню
type RefreshMenuRequest struct {
	RefreshNameAndDescription       bool `json:"refreshNameAndDescription"`
	RefreshPrice                    bool `json:"refreshPrice"`
	RefreshImages                   bool `json:"refreshImages"`
	RefreshModifiersNumber          bool `json:"refreshModifiersNumber"`
	RefreshNutritionPerHundredGrams bool `json:"refreshNutritionPerHundredGrams"`
	RefreshAllergens                bool `json:"refreshAllergens"`
	RefreshCombos                   bool `json:"refreshCombos"`
}
