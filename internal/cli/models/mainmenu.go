package models

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/api/inventory"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	guildStorage "lesta-start-battleship/cli/storage/guild"
	"strings"
)

type MainMenuModel struct {
	id       int
	username string
	gold     int
	selected int
	errorMsg string
	Clients  *clientdeps.Client
}

func NewMainMenuModel(id int, username string, gold int, clients *clientdeps.Client) *MainMenuModel {
	return &MainMenuModel{
		id:       id,
		username: username,
		gold:     gold,
		selected: 0,
		errorMsg: "",
		Clients:  clients,
	}
}

func (m *MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m *MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			m.selected = (m.selected - 1 + 6) % 6
			return m, nil

		case tea.KeyDown:
			m.selected = (m.selected + 1) % 6
			return m, nil

		case tea.KeyEnter:
			switch m.selected {
			case 0: // Бой
				return NewMatchmakingModel(m, m.id, m.username), nil
			case 1: // Инвентарь
				return m, m.loadHandler
			case 2: // Магазин
				model := NewShopModel(m, m.id, m.username, m.gold, ShopResponse{}, m.Clients)
				return model, model.Init()
			case 3: // Гильдия
				return m, m.guildHandler
			case 4: // Редактирование профиля
				return NewEditProfileModel(m.id, m.username, m.gold, m.Clients), nil
			case 5: // Рейтинг
				return NewScoreboardModel(m, m.id, m.username, m.gold, m.Clients), nil
			}
			return m, nil

		case tea.KeyEsc:
			return m, m.logoutHandler

		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case *inventory.UserInventoryResponse:
		return NewInventoryModel(m.id, m.username, m.gold, msg, m.Clients), nil

	case GuildDataMsg:
		return NewGuildModel(m.id, m.username, m.gold, msg.Member, msg.Guild, m.Clients), nil

	case GuildNoMemberMsg:
		return NewGuildModel(m.id, m.username, m.gold, nil, nil, m.Clients), nil
	}

	return m, nil
}

func (m *MainMenuModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Морской Бой"))
	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.username))
	sb.WriteString("\n\n")

	menuItems := []string{
		"⚔️  Бой",
		"🎒 Инвентарь",
		"🏪 Магазин",
		"🏰 Гильдия",
		"👤 Редактирование профиля",
		"🏆 Рейтинги",
	}

	for i, item := range menuItems {
		if i == m.selected {
			sb.WriteString(ui.SelectedStyle.Render("> " + item))
		} else {
			sb.WriteString(ui.NormalStyle.Render("  " + item))
		}
		sb.WriteString("\n")
	}

	if m.errorMsg != "" {
		sb.WriteString(ui.ErrorStyle.Render(m.errorMsg + "\n"))
	}

	sb.WriteString("\n")
	sb.WriteString(ui.HelpStyle.Render("↑/↓ - выбор, Enter - подтвердить, Esc - выход"))

	return sb.String()
}

func (m *MainMenuModel) loadHandler() tea.Msg {
	switch m.selected {
	case 1:
		ctx := context.Background()
		items, err := m.Clients.InventoryClient.GetUserInventory(ctx)
		if err != nil {
			m.errorMsg = fmt.Sprintf("%v", err)
			return m
		}
		return items
	}
	return nil
}

// Новый обработчик для гильдий
func (m *MainMenuModel) guildHandler() tea.Msg {
	ctx := context.Background()
	member, err := m.Clients.GuildsClient.GetMemberByUserID(ctx, m.id)
	if err != nil || member == nil {
		// Не состоит в гильдии
		return GuildNoMemberMsg{}
	}
	guildStorage.Self = *member
	guild, err := m.Clients.GuildsClient.GetGuildByTag(ctx, member.GuildTag)
	if err != nil || guild == nil {
		return GuildNoMemberMsg{}
	}
	return GuildDataMsg{
		Member: member,
		Guild:  guild,
	}
}

func (m *MainMenuModel) logoutHandler() tea.Msg {
	ctx := context.Background()
	err := m.Clients.AuthClient.Logout(ctx)
	if err != nil {
		m.errorMsg = err.Error()
		return m
	}
	return LogoutMsg{}
}
