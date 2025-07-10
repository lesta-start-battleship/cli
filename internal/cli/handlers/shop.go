package handlers

type ShopItem struct {
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
	Balance     int    `json:"balance"`
}

type ShopResponse struct {
	Balance int        `json:"balance"`
	Items   []ShopItem `json:"items"`
}

// ItemsHandler возвращает обычные предметы магазина
func ItemsHandler(token string) (ShopResponse, error) {
	return ShopResponse{
		Balance: 500,
		Items: []ShopItem{
			{
				Name:        "Крест Нахимова",
				Price:       150,
				Currency:    "gold",
				Description: "Позволяет разведать 5 клеток за раз",
			},
			{
				Name:        "Ремонтный набор",
				Price:       100,
				Currency:    "gold",
				Description: "Восстанавливает 1 клетку корабля",
			},
		},
	}, nil
}

// PromoHandler возвращает акционные товары
func PromoHandler(token string) (ShopResponse, error) {
	return ShopResponse{
		Balance: 500,
		Items: []ShopItem{
			{
				Name:        "Набор новичка (акция)",
				Price:       200,
				Currency:    "gold",
				Description: "2 креста Нахимова + 1 ремонтный набор (экономия 50 золота)",
			},
		},
	}, nil
}

// ChestsHandler возвращает сундуки
func ChestsHandler(token string) (ShopResponse, error) {
	return ShopResponse{
		Balance: 500,
		Items: []ShopItem{
			{
				Name:        "Гильдейский сундук",
				Price:       500,
				Currency:    "guild_rage",
				Description: "Содержит 3 случайных предмета для гильдейских войн",
			},
			{
				Name:        "Сундук удачи",
				Price:       300,
				Currency:    "gold",
				Description: "Содержит 1-3 случайных предмета",
			},
		},
	}, nil
}
