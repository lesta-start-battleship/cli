package shop

import "time"

// Product - игровой предмет для покупки
type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Currency    string `json:"currency_type"`
	Cost        int    `json:"cost"`
	DailyLimit  *int   `json:"daily_purchase_limit"`
	PromotionID *int   `json:"promotion"`
}

// Chest - игровой сундук
type Chest struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Gold            int    `json:"gold"`
	ItemProbability int    `json:"item_probability"`
	Experience      int    `json:"experience"`
	Currency        string `json:"currency_type"`
	Cost            int    `json:"cost"`
	DailyLimit      *int   `json:"daily_purchase_limit"`
	PromotionID     *int   `json:"promotion"`
}

// Promotion - активная акция
type Promotion struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	IsActive    bool      `json:"is_active"`
}

// Purchase - информация о покупке
type Purchase struct {
	ID          int       `json:"id"`
	UserID      int       `json:"owner"`
	Quantity    int       `json:"quantity"`
	Date        time.Time `json:"date"`
	ItemID      *int      `json:"item"`
	ChestID     *int      `json:"chest"`
	PromotionID *int      `json:"promotion"`
}

// OpenChestRequest - запрос на открытие сундука
type OpenChestRequest struct {
	ChestID int `json:"chest_id"`
	Amount  int `json:"amount"`
}
