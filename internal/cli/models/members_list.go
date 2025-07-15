package models

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/api/guilds"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	guildStorage "lesta-start-battleship/cli/storage/guild"
	"log"
	"strings"
)

const membersPerPage = 10

type MembersListModel struct {
	parent       tea.Model
	id           int
	username     string
	userRole     string
	guildTag     string
	guildName    string
	members      []guilds.MemberResponse
	currentPage  int
	totalPages   int
	selected     int
	actionMode   bool // true - выбор действия, false - выбор участников
	actionType   int  // 0 - изменить роль, 1 - удалить участника
	loading      bool
	errorMsg     string
	successMsg   string
	confirmState bool // true - подтверждение действия
	Clients      *clientdeps.Client
}

func NewMembersListModel(parent tea.Model, id int, username, userRole, guildTag, guildName string,
	clients *clientdeps.Client) *MembersListModel {
	return &MembersListModel{
		parent:      parent,
		id:          id,
		username:    username,
		userRole:    userRole,
		guildTag:    guildTag,
		guildName:   guildName,
		currentPage: 1,
		Clients:     clients,
	}
}

func (m *MembersListModel) Init() tea.Cmd {
	return m.loadMembers
}

func (m *MembersListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case MemberRoleChangeMsg:
		m.loading = false
		m.confirmState = false
		m.actionMode = false
		m.successMsg = fmt.Sprintf("Роль участника %s изменена", msg.Username)
		return m, m.loadMembers

	case MemberDeleteMsg:
		m.loading = false
		m.confirmState = false
		m.actionMode = false
		m.successMsg = fmt.Sprintf("Участник %s удален", msg.Username)
		return m, m.loadMembers
	}

	if m.confirmState {
		return m.handleConfirmState(msg)
	}

	if m.actionMode {
		return m.handleActionMode(msg)
	}

	return m.handleNormalMode(msg)
}

func (m *MembersListModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render(fmt.Sprintf("Участники гильдии %s [%s]", m.guildName, m.guildTag)))
	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render(fmt.Sprintf("Страница %d/%d", m.currentPage, m.totalPages)))
	sb.WriteString("\n\n")

	if m.loading {
		sb.WriteString("Загрузка списка участников...")
		return sb.String()
	}

	if m.errorMsg != "" {
		sb.WriteString(ui.ErrorStyle.Render(m.errorMsg))
		sb.WriteString("\n")
		return sb.String()
	}

	if m.confirmState {
		return m.renderConfirmView()
	}

	if m.actionMode {
		return m.renderActionView()
	}

	if len(m.members) == 0 {
		sb.WriteString("Нет участников в гильдии.\n")
	} else {
		for i, member := range m.members {
			line := fmt.Sprintf("%s (%s)", member.UserName, member.Role.Title)
			if i == m.selected {
				sb.WriteString(ui.SelectedStyle.Render(line))
			} else {
				sb.WriteString(ui.NormalStyle.Render(line))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	helpText := "↑/↓ - выбор, ←/→ - страницы, Esc - назад"
	if m.userRole == "owner" || m.userRole == "officer" {
		helpText += ", Enter - действия"
	}
	sb.WriteString(ui.HelpStyle.Render(helpText))

	return sb.String()
}

func (m *MembersListModel) renderActionView() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Выберите действие"))
	sb.WriteString("\n\n")

	// Только для owner
	if m.userRole == "owner" {
		changeRoleText := "Изменить роль"
		if m.actionType == 0 {
			changeRoleText = ui.SelectedStyle.Render("> " + changeRoleText)
		} else {
			changeRoleText = ui.NormalStyle.Render("  " + changeRoleText)
		}
		sb.WriteString(changeRoleText)
		sb.WriteString("\n")
	}

	// Удаление доступно и owner и officer
	deleteText := "Удалить из гильдии"
	if (m.userRole == "owner" && m.actionType == 1) || m.userRole == "officer" {
		deleteText = ui.SelectedStyle.Render("> " + deleteText)
	} else {
		deleteText = ui.NormalStyle.Render("  " + deleteText)
	}
	sb.WriteString(deleteText)
	sb.WriteString("\n")

	sb.WriteString("\n")
	helpText := "Enter - подтвердить, Esc - назад"
	if m.userRole == "owner" {
		helpText = "Tab - переключение, " + helpText
	}
	sb.WriteString(ui.HelpStyle.Render(helpText))

	return sb.String()
}

func (m *MembersListModel) renderConfirmView() string {
	var sb strings.Builder

	selectedMember := m.members[m.selected]

	if m.actionType == 0 {
		// Подтверждение изменения роли
		newRole := "officer"
		if selectedMember.Role.Title == "cabin_boy" {
			newRole = "officer"
		} else {
			newRole = "cabin_boy"
		}
		sb.WriteString(ui.TitleStyle.Render("Подтверждение"))
		sb.WriteString("\n\n")
		sb.WriteString(fmt.Sprintf("Изменить роль %s на %s?\n", selectedMember.UserName, newRole))
	} else {
		// Подтверждение удаления
		sb.WriteString(ui.TitleStyle.Render("Подтверждение"))
		sb.WriteString("\n\n")
		sb.WriteString(fmt.Sprintf("Вы точно хотите удалить %s из гильдии?\n", selectedMember.UserName))
	}

	sb.WriteString("\n")
	sb.WriteString(ui.HelpStyle.Render("Enter - подтвердить, Esc - отмена"))

	return sb.String()
}

func (m *MembersListModel) handleNormalMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.selected > 0 {
				m.selected--
			} else if m.currentPage > 1 {
				m.currentPage--
				m.selected = membersPerPage - 1
				return m, m.loadMembers
			}
			return m, nil

		case tea.KeyDown:
			if m.selected < len(m.members)-1 {
				m.selected++
			} else if m.currentPage < m.totalPages {
				m.currentPage++
				m.selected = 0
				return m, m.loadMembers
			}
			return m, nil

		case tea.KeyLeft:
			if m.currentPage > 1 {
				m.currentPage--
				m.selected = 0
				return m, m.loadMembers
			}
			return m, nil

		case tea.KeyRight:
			if m.currentPage < m.totalPages {
				m.currentPage++
				m.selected = 0
				return m, m.loadMembers
			}
			return m, nil

		case tea.KeyEnter:
			if len(m.members) > 0 && (m.userRole == "owner" || m.userRole == "officer") {
				m.actionMode = true
				m.actionType = 0
				return m, nil
			}
			return m, nil

		case tea.KeyEsc:
			return m.parent, nil
		}

	case *guilds.MemberPagination:
		m.loading = false
		m.members = msg.Items
		m.totalPages = msg.TotalPages
		if len(m.members) == 0 {
			m.errorMsg = "Участники не найдены"
		}
		for _, member := range m.members {
			guildStorage.SetMember(member.UserName, member)
		}
		return m, nil

	case error:
		m.loading = false
		m.errorMsg = msg.Error()
		return m, nil
	}
	return m, nil
}

func (m *MembersListModel) handleActionMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			if m.userRole == "owner" {
				m.actionType = (m.actionType + 1) % 2 // Переключение между действиями
			}
			return m, nil

		case tea.KeyEnter:
			m.confirmState = true
			return m, nil

		case tea.KeyEsc:
			m.actionMode = false
			return m, nil
		}
	}
	return m, nil
}

func (m *MembersListModel) handleConfirmState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.selected < 0 || m.selected >= len(m.members) {
				m.errorMsg = "Неверный выбор участника"
				m.confirmState = false
				return m, nil
			}

			selectedMember := m.members[m.selected]

			if m.actionType == 0 && m.userRole != "owner" {
				m.errorMsg = "Только владелец может изменять роли"
				m.confirmState = false
				return m, nil
			}

			m.loading = true
			m.errorMsg = ""
			m.successMsg = ""

			if m.actionType == 0 {
				var newRoleId int
				if selectedMember.Role.ID == 2 {
					newRoleId = 3
				} else {
					newRoleId = 2
				}

				return m, func() tea.Msg {
					ctx := context.Background()
					err := m.Clients.GuildsClient.EditMember(ctx, m.guildTag, selectedMember.UserID, m.id,
						guilds.EditMemberRequest{
							RoleID:   newRoleId,
							UserName: selectedMember.UserName,
						})
					if err != nil {
						log.Printf("Ошибка при изменении роли: %v", err.Error())
						return err
					}
					log.Printf("Новая роль участника %s: %d. Старая роль: %d", selectedMember.UserName, newRoleId, selectedMember.Role.ID)
					return MemberRoleChangeMsg{
						Username: selectedMember.UserName,
						NewRole:  getRole(newRoleId),
					}
				}
			} else {
				return m, func() tea.Msg {
					ctx := context.Background()
					err := m.Clients.GuildsClient.DeleteMember(ctx, m.guildTag, selectedMember.UserID, m.id)
					if err != nil {
						log.Printf("Ошибка при удалении: %v", err.Error())
						return err
					}
					return MemberDeleteMsg{Username: selectedMember.UserName}
				}
			}

		case tea.KeyEsc:
			m.confirmState = false
			return m, nil
		}
	}
	return m, nil
}

func getRole(roleId int) string {
	switch roleId {
	case 3:
		return "officer"
	default: // 2
		return "cabin_boy"
	}
}

func (m *MembersListModel) loadMembers() tea.Msg {
	m.loading = true
	ctx := context.Background()
	offset := (m.currentPage - 1) * membersPerPage
	members, err := m.Clients.GuildsClient.GetGuildMembers(ctx, m.guildTag, offset, membersPerPage)
	if err != nil {
		return err
	}
	return members
}
