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

const guildPerPage = 10

type GuildListModel struct {
	parent      tea.Model
	id          int
	username    string
	guilds      []guilds.GuildResponse
	currentPage int
	totalPages  int
	loading     bool
	errorMsg    string
	Clients     *clientdeps.Client
}

func NewGuildListModel(parent tea.Model, id int, username string, clients *clientdeps.Client) *GuildListModel {
	return &GuildListModel{
		parent:      parent,
		id:          id,
		username:    username,
		currentPage: 1,
		Clients:     clients,
	}
}

func (m *GuildListModel) Init() tea.Cmd {
	return m.loadGuilds
}

func (m *GuildListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			if m.currentPage > 1 {
				m.currentPage--
				return m, m.loadGuilds
			}
			return m, nil

		case tea.KeyRight:
			if m.currentPage < m.totalPages {
				m.currentPage++
				return m, m.loadGuilds
			}
			return m, nil

		case tea.KeyEsc:
			return m.parent, nil
		}

	case *guilds.GuildPagination:
		m.loading = false
		for _, guild := range msg.Items {
			guildStorage.SetGuild(guild.Tag, guild)
			guildStorage.SetGuildID(guild.ID, guild)
		}
		m.guilds = msg.Items
		m.totalPages = msg.TotalPages
		if len(m.guilds) == 0 {
			m.errorMsg = "Список гильдий пуст"
		}
		return m, nil

	case error:
		m.loading = false
		m.errorMsg = msg.Error()
		return m, nil
	}

	return m, nil
}

func (m *GuildListModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Список гильдий"))
	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render(fmt.Sprintf("Страница %d/%d", m.currentPage, m.totalPages)))
	sb.WriteString("\n\n")

	if m.loading {
		sb.WriteString(ui.NormalStyle.Render("Загрузка списка гильдий...\n"))
	}

	if m.errorMsg != "" {
		sb.WriteString(ui.ErrorStyle.Render(m.errorMsg + "\n"))
	}

	if len(m.guilds) == 0 {
		sb.WriteString(ui.NormalStyle.Render("Список гильдий пуст"))
	} else {
		for _, guild := range m.guilds {
			line := fmt.Sprintf("%s - %s [%s]", guild.Title, guild.Description, guild.Tag)

			if guild.IsFull {
				sb.WriteString(ui.ErrorStyle.Render(line + " (Полная)\n"))
			} else if !guild.IsActive {
				sb.WriteString(ui.WarningStyle.Render(line + " (Не участвует в войне)\n"))
			} else {
				sb.WriteString(ui.NormalStyle.Render(line + "\n"))
			}
		}
	}

	sb.WriteString("\n")
	sb.WriteString(ui.HelpStyle.Render("←/→ - переключение страниц, Esc - назад"))

	return sb.String()
}

func (m *GuildListModel) loadGuilds() tea.Msg {
	m.loading = true
	ctx := context.Background()
	offset := (m.currentPage - 1) * guildPerPage
	guildsList, err := m.Clients.GuildsClient.GetGuilds(ctx, offset, guildPerPage)
	log.Printf("Количество гильдий: %d", len(guildsList.Items))
	if err != nil {
		return err
	}
	return guildsList
}
