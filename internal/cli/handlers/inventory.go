package handlers

type InventoryItems struct {
	Name            string `json:"name"`
	Quantity        int    `json:"quantity"`
	ItemDescription string `json:"item_description"`
}

type InventoryResponse []InventoryItems

func InventoryHandler(token string) (InventoryResponse, error) {

	// Всегда возвращаем тестовые данные
	return InventoryResponse{
		{
			Name:            "Крест Нахимова",
			Quantity:        2,
			ItemDescription: "Один раз за игру позволяет просмотреть заданную клетку и клетки сверху, снизу, справа и слева от неё.",
		},
		{
			Name:            "Ремонтный набор",
			Quantity:        1,
			ItemDescription: "Позволяет починить одну клетку корабля, который полностью не потоплен.",
		},
	}, nil
}
