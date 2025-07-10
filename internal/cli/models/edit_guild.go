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

type EditGuildModel struct {
	parent       tea.Model
	id           int
	username     string
	gold         int
	guildTag     string
	originalName string
	originalDesc string
	name         string
	description  string
	activeField  int // 0 - name, 1 - description
	loading      bool
	errorMsg     string
	successMsg   string
	Clients      *clientdeps.Client
}

func NewEditGuildModel(parent tea.Model, id int, username string, gold int, guildTag, guildName, guildDesc string,
	clients *clientdeps.Client) *EditGuildModel {
	return &EditGuildModel{
		parent:       parent,
		id:           id,
		username:     username,
		gold:         gold,
		guildTag:     guildTag,
		originalName: guildName,
		originalDesc: guildDesc,
		name:         guildName,
		description:  guildDesc,
		Clients:      clients,
	}
}

func (m *EditGuildModel) Init() tea.Cmd {
	return nil
}

func (m *EditGuildModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.activeField == 0 {
				m.activeField = 1
				return m, nil
			}

			if m.name == m.originalName && m.description == m.originalDesc {
				return m.parent, nil
			}

			m.loading = true
			return m, func() tea.Msg {
				ctx := context.Background()
				guild, err := m.Clients.GuildsClient.EditGuild(ctx, m.guildTag, m.id, guilds.EditGuildRequest{
					Title:       m.name,
					Description: m.description,
				})
				if err != nil {
					return err
				}
				return GuildEditMsg{
					Guild: guild,
				}
			}

		case tea.KeyTab:
			m.activeField = (m.activeField + 1) % 2
			return m, nil

		case tea.KeyBackspace:
			switch m.activeField {
			case 0:
				if len(m.name) > 0 {
					m.name = m.name[:len(m.name)-1]
				}
			case 1:
				if len(m.description) > 0 {
					m.description = m.description[:len(m.description)-1]
				}
			}
			return m, nil

		case tea.KeyRunes:
			switch m.activeField {
			case 0:
				m.name += string(msg.Runes)
			case 1:
				m.description += string(msg.Runes)
			}
			return m, nil

		case tea.KeySpace:
			if m.activeField == 0 {
				m.name += " "
			} else {
				m.description += " "
			}
			return m, nil

		case tea.KeyEsc:
			return m.parent, nil
		}

	case GuildEditMsg:
		ctx := context.Background()
		member, err := m.Clients.GuildsClient.GetMemberByUserID(ctx, m.id)
		if err != nil {
			return m, func() tea.Msg {
				return err
			}
		}
		guildStorage.Self = *member
		return NewGuildModel(m.id, m.username, m.gold, member, msg.Guild, m.Clients), nil

	case error:
		m.loading = false
		m.errorMsg = msg.Error()
		return m, nil
	}

	return m, nil
}

func (m *EditGuildModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render(fmt.Sprintf("Изменение гильдии [%s]", m.guildTag)))
	sb.WriteString("\n\n")

	if m.loading {
		sb.WriteString("Сохранение изменений...")
		return sb.String()
	}

	if m.errorMsg != "" {
		sb.WriteString(ui.ErrorStyle.Render(m.errorMsg))
		sb.WriteString("\n\n")
	}

	// Название гильдии
	sb.WriteString("Название гильдии:\n")
	if m.activeField == 0 {
		sb.WriteString(ui.SelectedStyle.Render("> " + m.name + "_"))
	} else {
		sb.WriteString("  " + m.name)
	}
	sb.WriteString("\n\n")

	// Описание гильдии
	sb.WriteString("Описание гильдии:\n")
	if m.activeField == 1 {
		sb.WriteString(ui.SelectedStyle.Render("> " + m.description + "_"))
	} else {
		sb.WriteString("  " + m.description)
	}

	sb.WriteString("\n\n")
	sb.WriteString(ui.HelpStyle.Render("Tab - переключение полей, Enter - подтвердить, Esc - отмена"))

	return sb.String()
}
