package inventory

const (
	UserInventoryPath = "user_inventory"
)

// InventoryItem - элемент инвентаря
type InventoryItem struct {
	ItemID      int    `json:"item_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
}

// UserInventoryResponse - ответ с инвентарем пользователя
type UserInventoryResponse struct {
	UserID int             `json:"user_id"`
	Items  []InventoryItem `json:"items"`
}

// ErrorResponse - ответ об ошибке
type ErrorResponse struct {
	Error string `json:"error"`
}

// ToDo создание/обновление предметов
