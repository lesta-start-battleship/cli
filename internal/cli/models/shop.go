package models

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"strings"
)

// ShopItem унифицирует данные для отображения в магазине
type ShopItem struct {
	ID          int
	Name        string
	Description string
	Price       int
	Currency    string
	Type        string // "product", "chest", "promotion"
	PromotionID *int   // Для отображения акционной метки
}

// ShopResponse содержит данные магазина
type ShopResponse struct {
	Balance int
	Items   []ShopItem
}

type ShopModel struct {
	parent   tea.Model
	id       int
	username string
	gold     int
	items    ShopResponse
	selected int
	category int // 0-предметы, 1-акции, 2-сундуки
	err      error
	Clients  *clientdeps.Client
}

func NewShopModel(parent tea.Model, id int, username string, gold int, items ShopResponse, clients *clientdeps.Client) *ShopModel {
	return &ShopModel{
		parent:   parent,
		id:       id,
		username: username,
		gold:     gold,
		items:    items,
		selected: 0,
		category: 0,
		Clients:  clients,
	}
}

func (m *ShopModel) Init() tea.Cmd {
	return m.loadItems
}

func (m *ShopModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ShopResponse:
		m.items = msg
		return m, nil

	case error:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			m.category = (m.category - 1 + 3) % 3
			return m, m.loadItems

		case tea.KeyRight:
			m.category = (m.category + 1) % 3
			return m, m.loadItems

		case tea.KeyUp:
			if len(m.items.Items) > 0 {
				m.selected = (m.selected - 1 + len(m.items.Items)) % len(m.items.Items)
			}
			return m, nil

		case tea.KeyDown:
			if len(m.items.Items) > 0 {
				m.selected = (m.selected + 1) % len(m.items.Items)
			}
			return m, nil

		case tea.KeyEnter:
			if len(m.items.Items) > 0 {
				selectedItem := m.items.Items[m.selected]
				ctx := context.Background()

				/*
					// проверка баланса на клиенте
					if selectedItem.Price > m.items.Balance {
						m.err = fmt.Errorf("Недостаточно средств")
						return m, nil
					} */

				// логика покупки
				if selectedItem.Type == "product" {
					err := m.Clients.ShopClient.BuyProduct(ctx, selectedItem.ID)
					if err != nil {
						m.err = err
						return m, nil
					}
				} else if selectedItem.Type == "chest" {
					err := m.Clients.ShopClient.BuyChest(ctx, selectedItem.ID)
					if err != nil {
						m.err = err
						return m, nil
					}
				} else if selectedItem.Type == "promotion" {
					m.err = fmt.Errorf("Акции нельзя купить напрямую")
					return m, nil
				}

				// обновление баланса после покупки
				profile, err := m.Clients.AuthClient.GetProfile(ctx)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.items.Balance = profile.Currency.Gold
			}
			return m, nil

		case tea.KeyEsc:
			return NewMainMenuModel(m.id, m.username, m.gold, m.Clients), nil
		}
	}
	return m, nil
}

func (m *ShopModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Магазин"))
	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render(fmt.Sprintf("Пользователь: %s					Balance: %d 💰", m.username, m.gold)))
	sb.WriteString("\n\n")

	// Отображение категорий
	categories := []string{"Предметы", "Акции", "Сундуки"}
	for i, cat := range categories {
		if i == m.category {
			sb.WriteString(ui.SelectedStyle.Render("[" + cat + "] "))
		} else {
			sb.WriteString(ui.NormalStyle.Render(cat + " "))
		}
	}
	sb.WriteString("\n\n")

	if m.err != nil {
		sb.WriteString(ui.ErrorStyle.Render("Ошибка: " + m.err.Error()))
		return sb.String()
	}

	if len(m.items.Items) == 0 {
		sb.WriteString(ui.NormalStyle.Render("Товары отсутствуют"))
		return sb.String()
	}

	for i, item := range m.items.Items {
		if i == m.selected {
			sb.WriteString(ui.SelectedStyle.Render("> " + item.Name))
		} else {
			sb.WriteString(ui.NormalStyle.Render("  " + item.Name))
		}
		if item.Type != "promotion" {
			sb.WriteString(ui.NormalStyle.Render(fmt.Sprintf(" - %d %s", item.Price, item.Currency)))
		}
		if item.PromotionID != nil {
			sb.WriteString(ui.NormalStyle.Render(" [Акция]"))
		}
		sb.WriteString("\n")
		sb.WriteString(ui.NormalStyle.Render("   " + item.Description))
		sb.WriteString("\n\n")
	}

	sb.WriteString("\n")
	sb.WriteString(ui.HelpStyle.Render("←/→ - переключение категорий, ↑/↓ - выбор, Enter - купить, Esc - назад"))

	return sb.String()
}

func (m *ShopModel) loadItems() tea.Msg {
	ctx := context.Background()
	var items []ShopItem
	var err error

	switch m.category {
	case 0: // предметы
		products, err := m.Clients.ShopClient.GetProducts(ctx)
		if err != nil {
			return err
		}
		for _, p := range products {
			description := p.Description
			if p.PromotionID != nil {
				description += " (Акция)"
			}
			items = append(items, ShopItem{
				ID:          p.ID,
				Name:        p.Name,
				Description: description,
				Price:       p.Cost,
				Currency:    p.Currency,
				Type:        "product",
				PromotionID: p.PromotionID,
			})
		}
	case 1: // акции
		promotions, err := m.Clients.ShopClient.GetPromotions(ctx)
		if err != nil {
			return err
		}
		for _, p := range promotions {
			items = append(items, ShopItem{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       0,
				Currency:    "",
				Type:        "promotion",
				PromotionID: &p.ID,
			})
		}
	case 2: // сундуки
		chests, err := m.Clients.ShopClient.GetChests(ctx)
		if err != nil {
			return err
		}
		for _, c := range chests {
			description := fmt.Sprintf("Золото: %d, Вероятность предмета: %d%%, Опыт: %d", c.Gold, c.ItemProbability, c.Experience)
			if c.PromotionID != nil {
				description += " (Акция)"
			}
			items = append(items, ShopItem{
				ID:          c.ID,
				Name:        c.Name,
				Description: description,
				Price:       c.Cost,
				Currency:    c.Currency,
				Type:        "chest",
				PromotionID: c.PromotionID,
			})
		}
	}

	// получение баланса
	profile, err := m.Clients.AuthClient.GetProfile(ctx)
	if err != nil {
		return err
	}
	balance := profile.Currency.Gold

	return ShopResponse{
		Balance: balance,
		Items:   items,
	}
}
