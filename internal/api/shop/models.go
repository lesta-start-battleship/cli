package shop

// Product - игровой предмет для покупки
type Product struct {
	ID          int               `json:"item_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Currency    string            `json:"currency_type"`
	Cost        int               `json:"cost"`
	Promotion   *ProductPromotion `json:"promotion"`
	DailyLimit  *int              `json:"daily_purchase_limit"`
}

type ProductPromotion struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ChestProduct struct {
	ID          int    `json:"item_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Chest - игровой сундук
type Chest struct {
	ID              int            `json:"item_id"`
	Name            string         `json:"name"`
	Gold            int            `json:"gold"`
	PromotionID     *int           `json:"promotion"`
	ItemProbability int            `json:"item_probability"`
	Currency        string         `json:"currency_type"`
	Cost            int            `json:"cost"`
	Experience      int            `json:"experience"`
	Products        []ChestProduct `json:"products"`
	SpecialProducts []ChestProduct `json:"special_products"`
}

// Promotion - активная акция
type Promotion struct {
	ID        int                `json:"id"`
	Name      string             `json:"name"`
	StartDate string             `json:"start_date"`
	EndDate   string             `json:"end_date"`
	Duration  string             `json:"duration"`
	IsActive  string             `json:"is_active"`
	Chests    []Chest            `json:"chests"`
	Products  []ProductPromotion `json:"products"`
}

// Purchase - информация о покупке
type Purchase struct {
	ID          int    `json:"id"`
	UserID      int    `json:"owner"`
	Quantity    int    `json:"quantity"`
	Date        string `json:"date"`
	ItemID      *int   `json:"item"`
	ChestID     *int   `json:"chest"`
	PromotionID *int   `json:"promotion"`
}

// OpenChestRequest - запрос на открытие сундука
type OpenChestRequest struct {
	ChestID int `json:"item_id"`
	Amount  int `json:"amount"`
}

type ChestResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []Chest `json:"results"`
}
type ProductResponse struct {
	Count    int       `json:"count"`
	Next     *string   `json:"next"`
	Previous *string   `json:"previous"`
	Results  []Product `json:"results"`
}
type PromotionResponse struct {
	Count    int         `json:"count"`
	Next     *string     `json:"next"`
	Previous *string     `json:"previous"`
	Results  []Promotion `json:"results"`
}
