package models

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/api/inventory"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"strings"
)

type InventoryModel struct {
	id          int
	username    string
	gold        int
	items       *inventory.UserInventoryResponse
	selected    int
	showDetails bool
	Clients     *clientdeps.Client
}

func NewInventoryModel(id int, username string, gold int, items *inventory.UserInventoryResponse, clients *clientdeps.Client) *InventoryModel {
	return &InventoryModel{
		id:          id,
		username:    username,
		gold:        gold,
		items:       items,
		selected:    0,
		showDetails: false,
		Clients:     clients,
	}
}

func (m *InventoryModel) Init() tea.Cmd {
	return nil
}

func (m *InventoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if !m.showDetails && len(m.items.Items) > 0 {
				m.selected = (m.selected - 1 + len(m.items.Items)) % len(m.items.Items)
			}
			return m, nil

		case tea.KeyDown:
			if !m.showDetails && len(m.items.Items) > 0 {
				m.selected = (m.selected + 1) % len(m.items.Items)
			}
			return m, nil

		case tea.KeyEnter:
			if len(m.items.Items) > 0 && !m.showDetails {
				m.showDetails = true
			} else {
				m.showDetails = false
			}
			return m, nil

		case tea.KeyEsc:
			return NewMainMenuModel(m.id, m.username, m.gold, m.Clients), nil

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *InventoryModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Инвентарь"))
	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.username))
	sb.WriteString("\n\n")

	if len(m.items.Items) == 0 {
		sb.WriteString(ui.NormalStyle.Render("Инвентарь пуст"))
		return sb.String()
	}

	if m.showDetails {
		item := m.items.Items[m.selected]
		sb.WriteString(ui.SelectedStyle.Render(item.Name))
		sb.WriteString("\n\n")
		sb.WriteString(ui.NormalStyle.Render("Количество: "))
		sb.WriteString(ui.NormalStyle.Render(fmt.Sprintf("%d", item.Amount)))
		sb.WriteString("\n\n")
		sb.WriteString(ui.NormalStyle.Render("Описание:\n"))
		sb.WriteString(ui.NormalStyle.Render(item.Description))
	} else {
		for i, item := range m.items.Items {
			if i == m.selected {
				sb.WriteString(ui.SelectedStyle.Render("> " + item.Name))
				sb.WriteString(ui.NormalStyle.Render(fmt.Sprintf(" (x%d)", item.Amount)))
			} else {
				sb.WriteString(ui.NormalStyle.Render("  " + item.Name))
				sb.WriteString(ui.NormalStyle.Render(fmt.Sprintf(" (x%d)", item.Amount)))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	if m.showDetails {
		sb.WriteString(ui.HelpStyle.Render("Enter - назад, Esc - в меню"))
	} else {
		sb.WriteString(ui.HelpStyle.Render("↑/↓ - выбор, Enter - подробности, Esc - в меню"))
	}

	return sb.String()
}
