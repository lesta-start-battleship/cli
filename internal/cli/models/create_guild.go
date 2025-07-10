package models

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/api/guilds"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"strings"
)

type CreateGuildModel struct {
	parent      tea.Model
	id          int
	username    string
	gold        int
	name        string
	tag         string
	description string
	activeField int // 0 - name, 1 - tag, 2 - description
	loading     bool
	errorMsg    string
	Clients     *clientdeps.Client
}

func NewCreateGuildModel(parent tea.Model, id int, username string, gold int, clients *clientdeps.Client) *CreateGuildModel {
	return &CreateGuildModel{
		parent:   parent,
		id:       id,
		username: username,
		gold:     gold,
		Clients:  clients,
	}
}

func (m *CreateGuildModel) Init() tea.Cmd {
	return nil
}

func (m *CreateGuildModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.activeField < 2 {
				m.activeField++
				return m, nil
			}

			if len(m.name) < 3 {
				m.errorMsg = "Название гильдии должно быть не менее 3 символов."
				return m, nil
			}
			if len(m.tag) < 2 || len(m.tag) > 5 {
				m.errorMsg = "Тег гильдии должен быть от 2 до 5 символов."
				return m, nil
			}

			m.loading = true
			return m, func() tea.Msg {
				ctx := context.Background()
				guild, err := m.Clients.GuildsClient.CreateGuild(ctx, m.id, guilds.CreateGuildRequest{
					Title:       m.name,
					Tag:         m.tag,
					Description: m.description,
				})
				if err != nil {
					return err
				}
				return GuildCreatedMsg{Guild: guild}
			}

		case tea.KeyTab:
			m.activeField = (m.activeField + 1) % 3
			return m, nil

		case tea.KeyBackspace:
			switch m.activeField {
			case 0:
				if len(m.name) > 0 {
					m.name = m.name[:len(m.name)-1]
				}
			case 1:
				if len(m.tag) > 0 {
					m.tag = m.tag[:len(m.tag)-1]
				}
			case 2:
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
				m.tag += string(msg.Runes)
			case 2:
				m.description += string(msg.Runes)
			}
			return m, nil

		case tea.KeySpace:
			switch m.activeField {
			case 0:
				m.name += " "
			case 1:
				m.tag += ""
			case 2:
				m.description += ""
			}
			return m, nil

		case tea.KeyEsc:
			return m.parent, nil
		}

	case GuildCreatedMsg:
		ctx := context.Background()
		member, err := m.Clients.GuildsClient.GetMemberByUserID(ctx, m.id)
		if err != nil {
			m.errorMsg = err.Error()
			return m, nil
		}
		return NewGuildModel(m.id, m.username, m.gold, member, msg.Guild, m.Clients), nil

	case error:
		m.loading = false
		m.errorMsg = msg.Error()
		return m, nil
	}
	return m, nil
}

func (m *CreateGuildModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Создание гильдии"))
	sb.WriteString("\n\n")

	if m.loading {
		sb.WriteString("Cоздание гильдии...")
		return sb.String()
	}

	if m.errorMsg != "" {
		sb.WriteString(ui.ErrorStyle.Render(m.errorMsg))
		sb.WriteString("\n\n")
	}

	sb.WriteString("Название гильдии:\n")
	if m.activeField == 0 {
		sb.WriteString(ui.SelectedStyle.Render("> " + m.name + "_"))
	} else {
		sb.WriteString(" " + m.name)
	}
	sb.WriteString("\n\n")

	sb.WriteString("Тег гильдии:\n")
	if m.activeField == 1 {
		sb.WriteString(ui.SelectedStyle.Render("> " + m.tag + "_"))
	} else {
		sb.WriteString(" " + m.tag)
	}
	sb.WriteString("\n\n")

	sb.WriteString("Описание гильдии (необязательно):\n")
	if m.activeField == 2 {
		sb.WriteString(ui.SelectedStyle.Render("> " + m.description + "_"))
	} else {
		sb.WriteString(" " + m.description)
	}

	sb.WriteString("\n\n")
	sb.WriteString(ui.HelpStyle.Render("Enter - подтвердить, Tab - переключение полей, Esc - назад"))

	return sb.String()
}
