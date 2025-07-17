package models

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/api/guilds"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"strings"
)

const requestsPerPage = 10

type JoinRequestsModel struct {
	parent       tea.Model
	id           int
	username     string
	guildTag     string
	guildName    string
	requests     []guilds.RequestResponse
	currentPage  int
	totalPages   int
	selected     int
	selectedUser int
	confirmState bool
	loading      bool
	processing   bool
	errorMsg     string
	successMsg   string
	Clients      *clientdeps.Client
}

func NewJoinRequestsModel(parent tea.Model, id int, username string, guildTag string,
	guildName string, client *clientdeps.Client) *JoinRequestsModel {
	return &JoinRequestsModel{
		parent:      parent,
		id:          id,
		username:    username,
		guildTag:    guildTag,
		guildName:   guildName,
		currentPage: 1,
		loading:     true,
		Clients:     client,
	}
}

func (m *JoinRequestsModel) Init() tea.Cmd {
	return m.loadRequests
}

func (m *JoinRequestsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.confirmState {
		return m.handleConfirmState(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			m.errorMsg = ""
			m.successMsg = ""
			if m.selected > 0 {
				m.selected--
			} else if m.currentPage > 1 {
				m.currentPage--
				m.selected = requestsPerPage - 1
				return m, m.loadRequests
			}
			return m, nil

		case tea.KeyDown:
			m.errorMsg = ""
			m.successMsg = ""
			if m.selected < len(m.requests)-1 {
				m.selected++
			} else if m.currentPage < m.totalPages {
				m.currentPage++
				m.selected = 0
				return m, m.loadRequests
			}
			return m, nil

		case tea.KeyLeft:
			m.errorMsg = ""
			m.successMsg = ""
			if m.currentPage > 1 {
				m.currentPage--
				m.selected = 0
				return m, m.loadRequests
			}
			return m, nil

		case tea.KeyRight:
			m.errorMsg = ""
			m.successMsg = ""
			if m.currentPage < m.totalPages {
				m.currentPage++
				m.selected = 0
				return m, m.loadRequests
			}
			return m, nil

		case tea.KeyEnter:
			if len(m.requests) > 0 {
				m.selectedUser = m.requests[m.selected].UserID
				m.confirmState = true
			}
			return m, nil

		case tea.KeyEsc:
			return m.parent, nil
		}

	case *guilds.RequestPagination:
		m.loading = false
		m.requests = msg.Items
		m.totalPages = msg.TotalPages
		/*if len(m.requests) == 0 {
			m.errorMsg = "Заявки не найдены"
		}*/
		if m.selected >= len(m.requests) {
			m.selected = max(0, len(m.requests)-1)
		}
		return m, nil

	case RequestProcessedMsg:
		m.successMsg = msg.Message
		m.errorMsg = ""
		m.confirmState = false
		m.processing = false
		return m, nil

	case error:
		m.loading = false
		m.processing = false
		m.errorMsg = msg.Error()
		m.successMsg = ""
		return m, nil
	}

	return m, nil
}

func (m *JoinRequestsModel) handleConfirmState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Принять заявку
			m.confirmState = false
			m.processing = true
			return m, tea.Sequence(
				func() tea.Msg {
					ctx := context.Background()
					err := m.Clients.GuildsClient.ApplyJoinRequest(ctx, m.guildTag, m.selectedUser, m.id)
					if err != nil {
						return err
					}
					return RequestProcessedMsg{Message: "Заявка принята"}
				},
				m.loadRequests,
			)

		case tea.KeyBackspace:
			// Отклонить заявку
			m.confirmState = false
			m.processing = true
			return m, tea.Sequence(
				func() tea.Msg {
					ctx := context.Background()
					err := m.Clients.GuildsClient.CancelJoinRequest(ctx, m.guildTag, m.selectedUser, m.id)
					if err != nil {
						return err
					}
					return RequestProcessedMsg{Message: "Заявка отклонена"}
				},
				m.loadRequests,
			)

		case tea.KeyEsc:
			m.confirmState = false
			return m, nil
		}
	}

	return m, nil
}

func (m *JoinRequestsModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render(fmt.Sprintf("Заявки на вступление в гильдию %s [%s]", m.guildName, m.guildTag)))
	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render(fmt.Sprintf("Страница %d/%d", m.currentPage, m.totalPages)))
	sb.WriteString("\n\n")

	if m.loading {
		sb.WriteString("Загрузка списка заявок...")
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

	if m.confirmState {
		return m.renderConfirmView()
	}

	if len(m.requests) == 0 {
		sb.WriteString(ui.NormalStyle.Render("Нет заявок на вступление"))
	} else {
		for i, req := range m.requests {
			line := fmt.Sprintf("Usename: %s", req.UserName)
			//line := fmt.Sprintf("Username: %s, Gold: %d, Experience: %d, Rating: %d", req.Name, req.Gold, req.Experience, req.Rating)
			if i == m.selected {
				sb.WriteString(ui.SelectedStyle.Render("> " + line))
			} else {
				sb.WriteString(ui.NormalStyle.Render("  " + line))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	helpText := "↑/↓ - выбор, ←/→ - страницы, Enter - действия, Esc - назад"
	sb.WriteString(ui.HelpStyle.Render(helpText))

	return sb.String()
}

func (m *JoinRequestsModel) renderConfirmView() string {
	var sb strings.Builder

	selectedReq := m.requests[m.selected]
	sb.WriteString(ui.TitleStyle.Render("Обработка заявки"))
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("Игрок: %s\n\n", selectedReq.UserName))

	if m.processing {
		sb.WriteString(ui.HelpStyle.Render("Обработка запроса..."))
	} else {
		sb.WriteString(ui.HelpStyle.Render("Enter - принять заявку"))
		sb.WriteString("\n")
		sb.WriteString(ui.HelpStyle.Render("Backspace - отклонить заявку"))
	}
	sb.WriteString("\n")
	sb.WriteString(ui.HelpStyle.Render("Esc - отмена"))

	return sb.String()
}

func (m *JoinRequestsModel) loadRequests() tea.Msg {
	m.loading = true
	ctx := context.Background()
	reqID, err := m.Clients.GuildsClient.GetJoinRequests(ctx, m.guildTag, m.id)
	if err != nil {
		return err
	}

	/*var ids []int

	for _, user := range reqID.Items {
		ids = append(ids, user.UserID)
	}

	requests, err := m.Clients.ScoreboardClient.GetUserStats(ctx, ids, "", "",
		false, requestsPerPage, 1)*/

	return reqID
}
