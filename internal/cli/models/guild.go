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

type GuildModel struct {
	id            int
	username      string
	gold          int
	Member        *guilds.MemberResponse
	Guild         *guilds.GuildResponse
	selected      int
	inputGuildTag string // Тег гильдии
	deleteMode    bool   // Режим удаления гильдии
	isJoining     bool   // Режим вступления в гильдию
	isDeclareWar  bool   // Режим объявления войны
	loading       bool
	errorMsg      string
	successMsg    string
	Clients       *clientdeps.Client
}

func NewGuildModel(id int, username string, gold int, member *guilds.MemberResponse, guild *guilds.GuildResponse,
	clients *clientdeps.Client) *GuildModel {
	return &GuildModel{
		id:       id,
		username: username,
		gold:     gold,
		Member:   member,
		Guild:    guild,
		Clients:  clients,
	}
}

func (m *GuildModel) Init() tea.Cmd {
	return nil
}

func (m *GuildModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.isJoining {
		return m.handleJoinInput(msg)
	}

	if m.isDeclareWar {
		return m.handleDeclareWare(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}

		switch msg.Type {
		case tea.KeyUp:
			if m.deleteMode {
				return m, nil
			}
			menuItems := m.getMenuItems()
			m.selected = (m.selected - 1 + len(menuItems)) % len(menuItems)
			return m, nil

		case tea.KeyDown:
			if m.deleteMode {
				return m, nil
			}
			menuItems := m.getMenuItems()
			m.selected = (m.selected + 1) % len(menuItems)
			return m, nil

		case tea.KeyEnter:
			menuItems := m.getMenuItems()
			if len(menuItems) == 0 {
				return m, nil
			}

			selectedItem := menuItems[m.selected]

			if selectedItem == "Удалить гильдию" {
				if m.deleteMode {
					m.loading = true
					return m, func() tea.Msg {
						ctx := context.Background()
						err := m.Clients.GuildsClient.DeleteGuild(ctx, m.Guild.Tag, m.id)
						if err != nil {
							return err
						}
						return GuildDeletedMsg{}
					}
				} else {
					m.deleteMode = true
					return m, nil
				}
			} else {
				return m.handleMenuSelection()
			}

		case tea.KeyEsc:
			if m.deleteMode {
				m.deleteMode = false
				return m, nil
			}
			m.successMsg = ""
			m.errorMsg = ""
			return NewMainMenuModel(m.id, m.username, m.gold, m.Clients), nil
		}

	case GuildDeletedMsg:
		guildStorage.CleanStorage()
		m.successMsg = "Гильдия успешно удалена"
		return NewGuildModel(m.id, m.username, m.gold, nil, nil, m.Clients),
			func() tea.Msg {
				return m.successMsg
			}

	case JoinRequestSentMsg:
		m.loading = false
		m.isJoining = false
		m.successMsg = fmt.Sprintf("Запрос в гильдию [%s] отправлен", m.inputGuildTag)
		m.inputGuildTag = ""
		return m, nil

	case DeclareWarMsg:
		m.loading = false
		m.isDeclareWar = false
		m.successMsg = fmt.Sprintf("Война с гильдией [%s] объявлена", msg.TargetTag)
		m.inputGuildTag = ""
		return m, nil

	case error:
		m.loading = false
		m.errorMsg = fmt.Sprintf("Ошибка: %v", msg)
		return m, nil
	}

	return m, nil
}

func (m *GuildModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Гильдия"))
	sb.WriteString("\n")

	if m.loading {
		sb.WriteString("\nЗагрузка данных гильдии...")
		return sb.String()
	}

	if m.Member == nil || m.Guild == nil {
		sb.WriteString("\nВы не состоите в гильдии\n")
	} else {
		sb.WriteString(fmt.Sprintf("\nГильдия: [%s] %s\n", m.Guild.Tag, m.Guild.Title))
		sb.WriteString(fmt.Sprintf("Ваша роль: %s\n", m.Member.Role.Title))
	}

	sb.WriteString("\n")
	menuItems := m.getMenuItems()
	for i, item := range menuItems {
		line := " " + item

		if i == m.selected && item == "Вступить в гильдию" && m.isJoining {
			line = ui.SelectedStyle.Render("> Вступить в гильдию: " + m.inputGuildTag)
		} else if i == m.selected && item == "Объявить войну" && m.isDeclareWar {
			line = ui.SelectedStyle.Render("> Объявить войну: " + m.inputGuildTag)
		} else if i == m.selected {
			line = ui.SelectedStyle.Render("> " + item)
		}

		sb.WriteString(line)
		sb.WriteString("\n")
	}

	if m.deleteMode {
		sb.WriteString("\n")
		sb.WriteString(ui.ErrorStyle.Render("Вы уверены, что хотите удалить гильдию? Нажмите Enter для подтверждения"))
	}

	if m.successMsg != "" {
		sb.WriteString("\n")
		sb.WriteString(ui.SuccessStyle.Render(m.successMsg))
	}

	if m.errorMsg != "" {
		sb.WriteString(ui.ErrorStyle.Render(m.errorMsg + "\n"))
	}

	sb.WriteString("\n")
	helpText := "↑/↓ - выбор"
	if m.isJoining {
		helpText += ", Enter - отправить запрос, Esc - отмена"
	} else {
		helpText += ", Enter - подтвердить, Esc - назад"
	}
	sb.WriteString(ui.HelpStyle.Render(helpText))

	return sb.String()
}

func (m *GuildModel) getMenuItems() []string {
	if m.Member == nil || m.Guild == nil {
		return []string{"Вступить в гильдию", "Создать гильдию", "Список гильдий"}
	}

	switch m.Member.Role.Title {
	case "cabin_boy":
		return []string{"Список гильдий", "Список участников", "Чат гильдии", "Покинуть гильдию"}
	case "owner":
		return []string{"Объявить войну", "Запросы на войну", "Изменить гильдию", "Список участников", "Список гильдий",
			"Чат гильдии", "Запросы на вступление", "Удалить гильдию"}
	default:
		return []string{"Список гильдий", "Список участников", "Чат гильдии", "Запросы на вступление", "Покинуть гильдию"}
	}
}

func (m *GuildModel) handleMenuSelection() (tea.Model, tea.Cmd) {
	menuItems := m.getMenuItems()
	if len(menuItems) == 0 {
		return m, nil
	}

	selectedItem := menuItems[m.selected]

	switch selectedItem {
	case "Список гильдий":
		model := NewGuildListModel(m, m.id, m.username, m.Clients)
		return model, model.Init()
	case "Список участников":
		model := NewMembersListModel(m, m.id, m.username, m.Member.Role.Title, m.Guild.Tag, m.Guild.Title, m.Clients)
		return model, model.Init()
	case "Чат гильдии":
		// Инициализация чата гильдии с правильным guildID
		guildID := 0
		if m.Guild != nil {
			guildID = m.Guild.ID
		}
		return m, func() tea.Msg {
			return OpenChatMsg{
				GuildID: guildID,
			}
		}
	case "Покинуть гильдию":
		return NewExitGuildModel(m, m.id, m.username, m.gold, m.Guild.Tag, m.Guild.Title, m.Clients), nil
	case "Объявить войну":
		m.isDeclareWar = true
		m.inputGuildTag = ""
		return m, nil
	case "Запросы на войну":
		model := NewWarRequestsModel(m, m.id, m.username, m.Guild.Tag, m.Guild.Title, m.Guild.ID, m.Clients)
		return model, model.Init()
	case "Изменить гильдию":
		return NewEditGuildModel(m, m.id, m.username, m.gold, m.Guild.Tag, m.Guild.Title, m.Guild.Description, m.Clients), nil
	case "Запросы на вступление":
		model := NewJoinRequestsModel(m, m.id, m.username, m.Guild.Tag, m.Guild.Title, m.Clients)
		return model, model.Init()
	case "Создать гильдию":
		return NewCreateGuildModel(m, m.id, m.username, m.gold, m.Clients), nil
	case "Вступить в гильдию":
		m.isJoining = true
		m.inputGuildTag = ""
		return m, nil
	default:
		return m, nil
	}
}

func (m *GuildModel) handleJoinInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.inputGuildTag == "" {
				m.errorMsg = "Введите тег гильдия"
				return m, nil
			}
			m.loading = true
			m.errorMsg = ""
			m.isJoining = false
			return m, func() tea.Msg {
				ctx := context.Background()
				err := m.Clients.GuildsClient.SendJoinRequest(ctx, m.inputGuildTag, m.id)
				if err != nil {
					return err
				}
				return JoinRequestSentMsg{}
			}

		case tea.KeyBackspace:
			if len(m.inputGuildTag) > 0 {
				m.inputGuildTag = m.inputGuildTag[:len(m.inputGuildTag)-1]
			}
			return m, nil

		case tea.KeyRunes:
			m.inputGuildTag += string(msg.Runes)
			return m, nil

		case tea.KeyEsc:
			m.isJoining = false
			m.inputGuildTag = ""
			m.errorMsg = ""
			return m, nil
		}
	}
	return m, nil
}

func (m *GuildModel) handleDeclareWare(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.inputGuildTag == "" {
				m.errorMsg = "Введите тег гильдии"
				return m, nil
			}

			m.loading = true
			m.errorMsg = ""
			m.successMsg = ""
			m.isDeclareWar = false

			return m, func() tea.Msg {
				targetGuild, ok := guildStorage.GetGuild(m.inputGuildTag)
				if !ok {
					m.loading = false
					m.errorMsg = fmt.Sprintf("Гильдия с тегом [%s] не найдена", m.inputGuildTag)
					return nil
				}
				ctx := context.Background()
				_, err := m.Clients.GuildsClient.DeclareWar(ctx, m.Guild.ID, targetGuild.ID, m.id)
				if err != nil {
					return err
				}
				return DeclareWarMsg{
					TargetTag: m.inputGuildTag,
				}
			}

		case tea.KeyBackspace:
			if len(m.inputGuildTag) > 0 {
				m.inputGuildTag = m.inputGuildTag[:len(m.inputGuildTag)-1]
			}
			return m, nil

		case tea.KeyRunes:
			m.inputGuildTag += string(msg.Runes)
			return m, nil

		case tea.KeyEsc:
			m.isDeclareWar = false
			m.inputGuildTag = ""
			m.errorMsg = ""
			return m, nil
		}
	}
	return m, nil
}
