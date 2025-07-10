package models

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/api/guilds"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	guildStorage "lesta-start-battleship/cli/storage/guild"
	"strings"
)

const warsPerPage = 5

type WarRequestsModel struct {
	parent      tea.Model
	id          int
	username    string
	gold        int
	guildTag    string
	guildName   string
	guildID     int
	wars        []guilds.GuildWarItem
	currentPage int
	totalPages  int
	selected    int
	actionMode  bool // true - выбор действия
	loading     bool
	errorMsg    string
	successMsg  string
	Clients     *clientdeps.Client
}

func NewWarRequestsModel(parent tea.Model, id int, username string, guildTag, guildName string, guildID int, clients *clientdeps.Client) *WarRequestsModel {
	return &WarRequestsModel{
		parent:      parent,
		id:          id,
		username:    username,
		guildTag:    guildTag,
		guildName:   guildName,
		guildID:     guildID,
		currentPage: 1,
		Clients:     clients,
	}
}

func (m *WarRequestsModel) Init() tea.Cmd {
	return m.loadWarRequests
}

func (m *WarRequestsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.actionMode {
		return m.handleActionMode(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.selected > 0 {
				m.selected--
			} else if m.currentPage > 1 {
				m.currentPage--
				m.selected = warsPerPage - 1
				return m, m.loadWarRequests
			}
			return m, nil

		case tea.KeyDown:
			if m.selected < len(m.wars)-1 {
				m.selected++
			} else if m.currentPage < m.totalPages {
				m.currentPage++
				m.selected = 0
				return m, m.loadWarRequests
			}
			return m, nil

		case tea.KeyLeft:
			if m.currentPage > 1 {
				m.currentPage--
				m.selected = 0
				return m, m.loadWarRequests
			}
			return m, nil

		case tea.KeyRight:
			if m.currentPage < m.totalPages {
				m.currentPage++
				m.selected = 0
				return m, m.loadWarRequests
			}
			return m, nil

		case tea.KeyEnter:
			if len(m.wars) > 0 {
				m.actionMode = true
			}
			return m, nil

		case tea.KeyEsc:
			return m.parent, nil
		}

	case *guilds.GuildWarListResponse:
		m.loading = false
		m.wars = msg.Results
		m.totalPages = msg.TotalPages
		if len(m.wars) == 0 {
			m.errorMsg = "Запросы на войну не найдены"
		}
		return m, nil

	case WarRequestProcessedMsg:
		m.successMsg = msg.Message
		m.actionMode = false
		return m, m.loadWarRequests

	case error:
		m.loading = false
		m.errorMsg = msg.Error()
		return m, nil
	}

	return m, nil
}

func (m *WarRequestsModel) handleActionMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Принять войну
			selectedWar := m.wars[m.selected]
			return m, func() tea.Msg {
				ctx := context.Background()
				_, err := m.Clients.GuildsClient.ConfirmWar(ctx, selectedWar.ID, m.id)
				if err != nil {
					return err
				}
				return WarRequestProcessedMsg{
					Message: fmt.Sprintf("Война с гильдией ID %d подтверждена", selectedWar.InitiatorGuildID),
				}
			}

		case tea.KeyBackspace:
			// Отклонить войну
			selectedWar := m.wars[m.selected]
			return m, func() tea.Msg {
				ctx := context.Background()
				_, err := m.Clients.GuildsClient.CancelWar(ctx, selectedWar.ID, m.id)
				if err != nil {
					return err
				}
				return WarRequestProcessedMsg{
					Message: fmt.Sprintf("Война с гильдией ID %d отклонена", selectedWar.InitiatorGuildID),
				}
			}

		case tea.KeyEsc:
			m.actionMode = false
			return m, nil
		}
	}

	return m, nil
}

func (m *WarRequestsModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render(fmt.Sprintf("Запросы на войну %s [%s]", m.guildName, m.guildTag)))
	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render(fmt.Sprintf("Страница %d/%d", m.currentPage, m.totalPages)))
	sb.WriteString("\n\n")

	if m.loading {
		sb.WriteString("Загрузка списка запросов...")
		return sb.String()
	}

	if m.errorMsg != "" {
		sb.WriteString(ui.ErrorStyle.Render(m.errorMsg))
		sb.WriteString("\n\n")
	}

	if m.successMsg != "" {
		sb.WriteString(ui.SuccessStyle.Render(m.successMsg))
		sb.WriteString("\n\n")
	}

	if m.actionMode {
		return m.renderActionView()
	}

	var line string
	if len(m.wars) == 0 {
		sb.WriteString(ui.NormalStyle.Render("Нет запросов на войну"))
	} else {
		for i, war := range m.wars {
			initiatorGuild, ok := guildStorage.GetGuildID(war.InitiatorGuildID)
			if !ok {
				line = fmt.Sprintf("Гильдия ID: %d (Статус: %s)", war.InitiatorGuildID, war.Status)
			}
			line = fmt.Sprintf("Гильдия Tag: %s (Статус: %s)", initiatorGuild.Tag, war.Status)
			if i == m.selected {
				sb.WriteString(ui.SelectedStyle.Render("> " + line))
			} else {
				sb.WriteString("  " + line)
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	helpText := "↑/↓ - выбор, ←/→ - страницы, Enter - действия, Esc - назад"
	sb.WriteString(ui.NormalStyle.Render(helpText))

	return sb.String()
}

func (m *WarRequestsModel) renderActionView() string {
	var sb strings.Builder

	selectedWar := m.wars[m.selected]
	sb.WriteString(ui.TitleStyle.Render("Обработка запроса на войну"))
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("Гильдия ID: %d\n", selectedWar.InitiatorGuildID))
	sb.WriteString(fmt.Sprintf("Статус: %s\n\n", selectedWar.Status))
	sb.WriteString(ui.NormalStyle.Render("Enter - принять войну"))
	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render("Backspace - отклонить войну"))
	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render("Esc - отмена"))

	return sb.String()
}

func (m *WarRequestsModel) loadWarRequests() tea.Msg {
	m.loading = true
	ctx := context.Background()
	isInitiator := false
	isTarget := true
	status := guilds.WarStatusPending
	wars, err := m.Clients.GuildsClient.GetGuildWarList(
		ctx,
		m.id,
		m.guildID,
		&isInitiator, // isInitiator
		&isTarget,    // isTarget
		&status,      // только ожидающие подтверждения
		m.currentPage,
		warsPerPage,
	)
	if err != nil {
		return err
	}
	return wars
}
