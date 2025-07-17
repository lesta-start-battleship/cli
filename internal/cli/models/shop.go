package models

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/api/shop"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"strings"
	"time"
)

// ShopItem унифицирует данные для отображения в магазине
type ShopItem struct {
	ID          int
	Name        string
	Description string
	Price       int
	Currency    string
	Type        string // "product", "chest", "promotion"
	DailyLimit  *int
	Promotion   *PromotionInfo
	Gold        int
	Exp         int
	ItemProb    int
	StartDate   string
	EndDate     string
}

// ShopResponse содержит данные магазина
type PromotionInfo struct {
	ID   int
	Name string
}

type ShopModel struct {
	id         int
	username   string
	gold       int
	items      []ShopItem
	selected   int
	category   int // 0-предметы, 1-акции, 2-сундуки
	err        error
	loading    bool
	success    string
	Clients    *clientdeps.Client
	balance    map[string]int
	pageSize   int
	page       int
	totalPages int
}

func NewShopModel(id int, username string, gold int, clients *clientdeps.Client) *ShopModel {
	return &ShopModel{
		id:       id,
		username: username,
		gold:     gold,
		Clients:  clients,
		balance:  make(map[string]int),
		pageSize: 5,
	}
}

func (m *ShopModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadBalance,
		m.loadItems,
	)
}

func (m *ShopModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	/*if m.loading {
		return m, nil
	}*/

	switch msg := msg.(type) {
	case []ShopItem:
		m.items = msg
		m.loading = false
		m.totalPages = (len(m.items) + m.pageSize - 1) / m.pageSize
		return m, nil

	case map[string]int:
		m.balance = msg
		return m, nil

	case error:
		m.err = msg
		m.loading = false
		return m, nil

	case string:
		m.success = msg
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			if m.loading {
				return m, nil
			}
			m.err = nil
			m.category = (m.category - 1 + 3) % 3
			m.page = 0
			return m, m.loadItems

		case tea.KeyRight:
			if m.loading {
				return m, nil
			}
			m.err = nil
			m.category = (m.category + 1) % 3
			m.page = 0
			return m, m.loadItems

		case tea.KeyUp:
			m.err = nil
			if len(m.items) > 0 {
				m.selected = (m.selected - 1 + len(m.items)) % len(m.items)
			}
			return m, nil

		case tea.KeyDown:
			m.err = nil
			if len(m.items) > 0 {
				m.selected = (m.selected + 1) % len(m.items)
			}
			return m, nil

		case tea.KeyEnter:
			if m.loading || len(m.items) == 0 {
				return m, nil
			}
			return m, m.buyItem(m.items[m.selected])

		case tea.KeyEsc:
			return NewMainMenuModel(m.id, m.username, m.gold, m.Clients), nil
		}
	}

	return m, nil
}

func (m *ShopModel) View() string {
	var sb strings.Builder

	// Заголовок
	sb.WriteString(ui.TitleStyle.Render("🏪 Магазин"))
	sb.WriteString("\n\n")

	// Баланс
	sb.WriteString(m.renderBalance())
	sb.WriteString("\n")

	// Категории
	sb.WriteString(m.renderCategories())
	sb.WriteString("\n\n")

	if m.err != nil {
		sb.WriteString(ui.ErrorStyle.Render("Ошибка: " + m.err.Error() + "\n\n"))
	}
	if m.success != "" {
		sb.WriteString(ui.SuccessStyle.Render(m.success + "\n\n"))
	}

	// Товары
	if m.loading {
		sb.WriteString(ui.NormalStyle.Render("Загрузка товаров...\n"))
	} else if len(m.items) == 0 {
		sb.WriteString(ui.NormalStyle.Render("Товары отсутствуют\n"))
	} else {
		sb.WriteString(m.renderItems())
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("Страница %d/%d\n", m.page+1, m.totalPages))
	}

	sb.WriteString("\n")
	sb.WriteString(ui.HelpStyle.Render(
		"←/→ - категории | ↑/↓ - выбор | Enter - купить | Esc - назад",
	))

	return sb.String()
}

func (m *ShopModel) renderBalance() string {
	var balances []string
	for currency, amount := range m.balance {
		switch currency {
		case "gold":
			balances = append(balances, fmt.Sprintf("💰 Золото: %d", amount))
		case "guild_rage":
			balances = append(balances, fmt.Sprintf("⚡ Ярость гильдии: %d", amount))
		}
	}
	return ui.SubtitleStyle.Render(strings.Join(balances, " | "))
}

func (m *ShopModel) renderCategories() string {
	categories := []string{"Предметы", "Акции", "Сундуки"}
	var rendered []string
	for i, cat := range categories {
		if i == m.category {
			rendered = append(rendered, ui.SelectedTabStyle.Render(cat))
		} else {
			rendered = append(rendered, ui.NormalTabStyle.Render(cat))
		}
	}
	return strings.Join(rendered, " ")
}

func (m *ShopModel) renderItems() string {
	var sb strings.Builder
	start := m.page * m.pageSize
	end := start + m.pageSize
	if end > len(m.items) {
		end = len(m.items)
	}

	for i, item := range m.items[start:end] {
		// Выделение выбранного элемента
		prefix := "  "
		if i+start == m.selected {
			prefix = ui.SelectedStyle.Render("> ")
		}

		// Название и цена
		name := prefix + item.Name
		if item.Price > 0 {
			name += fmt.Sprintf(" - %d %s", item.Price, getCurrencySymbol(item.Currency))
		}

		// Акция
		if item.Promotion != nil {
			name += " " + ui.PromotionStyle.Render("[АКЦИЯ]")
		}

		sb.WriteString(name + "\n")

		// Описание
		desc := "   " + item.Description
		switch item.Type {
		case "chest":
			desc += fmt.Sprintf("\n   🎁 Содержимое: %d золота, %d опыта, %d%% шанс предмета",
				item.Gold, item.Exp, item.ItemProb)
		case "promotion":
			desc += fmt.Sprintf("\n   🕒 Срок: %s - %s",
				item.StartDate, item.EndDate)
		}

		if item.DailyLimit != nil {
			desc += fmt.Sprintf("\n   🛒 Лимит: %d/день", *item.DailyLimit)
		}

		sb.WriteString(desc + "\n\n")
	}

	return sb.String()
}

func getCurrencySymbol(currency string) string {
	switch currency {
	case "gold":
		return "💰"
	case "guild_rage":
		return "⚡"
	default:
		return currency
	}
}

func (m *ShopModel) loadBalance() tea.Msg {
	ctx := context.Background()
	profile, err := m.Clients.AuthClient.GetProfile(ctx)
	if err != nil {
		return err
	}

	balance := map[string]int{
		"gold":       profile.Currency.Gold,
		"guild_rage": profile.Currency.GuildRage,
	}

	return balance
}

func (m *ShopModel) loadItems() tea.Msg {
	m.loading = true
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var items []ShopItem

	switch m.category {
	case 0: // Предметы
		products, err := m.Clients.ShopClient.GetProducts(ctx)
		if err != nil {
			return err
		}
		for _, p := range products {
			items = append(items, ShopItem{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Cost,
				Currency:    p.Currency,
				Type:        "product",
				Promotion:   toPromotionInfo(p.Promotion),
				DailyLimit:  p.DailyLimit,
			})
		}

	case 1: // Акции
		promotions, err := m.Clients.ShopClient.GetPromotions(ctx)
		if err != nil {
			return err
		}
		for _, p := range promotions {
			items = append(items, ShopItem{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.IsActive,
				Type:        "promotion",
				StartDate:   p.StartDate,
				EndDate:     p.EndDate,
			})
		}

	case 2: // Сундуки
		chests, err := m.Clients.ShopClient.GetChests(ctx)
		if err != nil {
			return err
		}
		for _, c := range chests {
			items = append(items, ShopItem{
				ID:       c.ID,
				Name:     c.Name,
				Price:    c.Cost,
				Currency: c.Currency,
				Type:     "chest",
				Gold:     c.Gold,
				Exp:      c.Experience,
				ItemProb: c.ItemProbability,
				Promotion: &PromotionInfo{
					ID:   *c.PromotionID,
					Name: "Специальное предложение",
				},
			})
		}
	}

	return items
}

func (m *ShopModel) buyItem(item ShopItem) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		var err error

		switch item.Type {
		case "product":
			err = m.Clients.ShopClient.BuyProduct(ctx, item.ID)
		case "chest":
			err = m.Clients.ShopClient.BuyChest(ctx, item.ID)
		default:
			return fmt.Errorf("этот тип товара нельзя купить напрямую")
		}

		if err != nil {
			return err
		}

		// Обновляем баланс после покупки
		balance, ok := m.loadBalance().(map[string]int)
		if !ok {
			return fmt.Errorf("ошибка обновления баланса")
		}

		return tea.Batch(
			func() tea.Msg { return balance },
			func() tea.Msg { return fmt.Sprintf("Успешно куплено: %s", item.Name) },
			m.loadItems, // Перезагружаем список товаров
		)
	}
}

func toPromotionInfo(p *shop.ProductPromotion) *PromotionInfo {
	if p == nil {
		return nil
	}
	return &PromotionInfo{
		ID:   p.ID,
		Name: p.Name,
	}
}
