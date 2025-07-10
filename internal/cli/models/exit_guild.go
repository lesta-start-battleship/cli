package models

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	guildStore "lesta-start-battleship/cli/storage/guild"
	"strings"
)

type ExitGuildModel struct {
	parent      tea.Model
	id          int
	username    string
	gold        int
	guildTag    string
	guildName   string
	confirmStep bool
	loading     bool
	errorMsg    string
	successMsg  string
	Client      *clientdeps.Client
}

func NewExitGuildModel(parent tea.Model, id int, username string, gold int, guildTag string, guildName string,
	client *clientdeps.Client) *ExitGuildModel {
	return &ExitGuildModel{
		parent:    parent,
		id:        id,
		username:  username,
		gold:      gold,
		guildTag:  guildTag,
		guildName: guildName,
		Client:    client,
	}
}

func (m *ExitGuildModel) Init() tea.Cmd {
	return nil
}

func (m *ExitGuildModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if !m.confirmStep {
				m.confirmStep = true
				return m, nil
			}

			m.loading = true
			return m, func() tea.Msg {
				ctx := context.Background()
				err := m.Client.GuildsClient.ExitGuild(ctx, m.guildTag)
				if err != nil {
					return err
				}
				return GuildExitedMsg{}
			}

		case tea.KeyEsc:
			if m.confirmStep {
				m.confirmStep = false
				return m, nil
			} else {
				return m.parent, nil
			}
		}

	case GuildExitedMsg:
		guildStore.CleanStorage()
		return NewGuildModel(m.id, m.username, m.gold, nil, nil, m.Client), nil

	case error:
		m.loading = false
		m.errorMsg = msg.Error()
		return m, nil
	}

	return m, nil
}

func (m *ExitGuildModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Покинуть гильдию"))
	sb.WriteString("\n\n")

	if m.loading {
		sb.WriteString("Покидаем гильдию...")
		return sb.String()
	}

	if m.errorMsg != "" {
		sb.WriteString(ui.ErrorStyle.Render(m.errorMsg))
		sb.WriteString("\n\n")
	}

	if m.successMsg != "" {
		sb.WriteString(ui.SuccessStyle.Render(m.successMsg))
		sb.WriteString("\n\n")
		return sb.String()
	}

	if !m.confirmStep {
		sb.WriteString(fmt.Sprintf("Вы собираетесь покинуть гильдию [%s] %s\n\n", m.guildTag, m.guildName))
		sb.WriteString(ui.NormalStyle.Render("Enter - продолжить, Esc - отмена"))
	} else {
		sb.WriteString(ui.ErrorStyle.Render("Вы уверены, что хотите покинуть гильдию?"))
		sb.WriteString("\n\n")
		sb.WriteString(fmt.Sprintf("Гильдия: [%s] %s\n\n", m.guildTag, m.guildName))
		sb.WriteString(ui.HelpStyle.Render("Enter - подтвердить, Esc - отмена"))
	}

	return sb.String()
}
