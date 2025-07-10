package models

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/api/scoreboard"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"strings"
)

const (
	pageSize = 10
)

type ScoreboardModel struct {
	parent       tea.Model
	id           int
	username     string
	gold         int
	activeTab    int // 0-моя, 1-игроки, 2-гильдия
	playersTab   int // 0-золото, 1-опыт, 2-рейтинг, 3-сундуки
	guildsTab    int // 0-игроки, 1-победы
	myStats      *scoreboard.UserStat
	playersStats *scoreboard.UserListResponse
	guildStats   *scoreboard.GuildListResponse
	err          error
	currentPage  int
	totalPages   int
	tableWidth   int
	Clients      *clientdeps.Client
}

func NewScoreboardModel(parent tea.Model, id int, username string, gold int, clients *clientdeps.Client) *ScoreboardModel {
	return &ScoreboardModel{
		parent:      parent,
		id:          id,
		username:    username,
		gold:        gold,
		activeTab:   0,
		playersTab:  0,
		guildsTab:   0,
		currentPage: 1,
		tableWidth:  80,
		Clients:     clients,
	}
}

func (m *ScoreboardModel) Init() tea.Cmd {
	return m.loadStats
}

func (m *ScoreboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.tableWidth = msg.Width - 10
		return m, nil

	case *scoreboard.UserStat:
		m.myStats = msg
		return m, nil

	case *scoreboard.UserListResponse:
		m.playersStats = msg
		return m, nil

	case *scoreboard.GuildListResponse:
		m.guildStats = msg
		return m, nil

	case error:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			if m.activeTab == 0 {
				return m, nil
			}
			m.activeTab--
			m.currentPage = 1
			return m, m.loadStats

		case tea.KeyRight:
			if m.activeTab == 2 {
				return m, nil
			}
			m.activeTab++
			m.currentPage = 1
			return m, m.loadStats

		case tea.KeyDown:
			if m.activeTab == 1 && m.playersStats != nil && m.currentPage < m.playersStats.PageAmount {
				m.currentPage++
				return m, m.loadStats
			} else if m.activeTab == 2 && m.guildStats != nil && m.currentPage < m.guildStats.PageAmount {
				m.currentPage++
				return m, m.loadStats
			}
			return m, nil

		case tea.KeyUp:
			if m.currentPage > 1 {
				m.currentPage--
				return m, m.loadStats
			}
			return m, nil

		case tea.KeyTab:
			if m.activeTab == 1 {
				m.playersTab = (m.playersTab + 1) % 4
				m.currentPage = 1
				return m, m.loadStats
			} else if m.activeTab == 2 {
				m.guildsTab = (m.guildsTab + 1) % 2
				m.currentPage = 1
				return m, m.loadStats
			}
			return m, nil

		case tea.KeyEsc:
			return NewMainMenuModel(m.id, m.username, m.gold, m.Clients), nil
		}
	}
	return m, nil
}

func (m *ScoreboardModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Рейтинги"))
	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.username))
	sb.WriteString("\n\n")

	tabs := []string{"Моя статистика", "Игроки", "Моя гильдия"}
	for i, tab := range tabs {
		if i == m.activeTab {
			sb.WriteString(ui.SelectedStyle.Render(" [" + tab + "] "))
		} else {
			sb.WriteString(ui.NormalStyle.Render(tab + " "))
		}
	}
	sb.WriteString("\n\n")

	if m.err != nil {
		sb.WriteString(ui.ErrorStyle.Render("Ошибка: " + m.err.Error()))
		return sb.String()
	}

	switch m.activeTab {
	case 0:
		if m.myStats == nil {
			sb.WriteString(ui.NormalStyle.Render("Загрузка данных..."))
			break
		}
		sb.WriteString(fmt.Sprintf("Игрок: %s\n", m.myStats.Name))
		sb.WriteString(fmt.Sprintf("Золото: %d (позиция: %d)\n", m.myStats.Gold, m.myStats.GoldRatingPos))
		sb.WriteString(fmt.Sprintf("Опыт: %d (позиция: %d)\n", m.myStats.Experience, m.myStats.ExpRatingPos))
		sb.WriteString(fmt.Sprintf("Рейтинг: %d (позиция: %d)\n", m.myStats.Rating, m.myStats.RatingRatingPos))
		sb.WriteString(fmt.Sprintf("Открыто сундуков: %d (позиция: %d)\n", m.myStats.ChestsOpened, m.myStats.ChestsOpenedRatingPos))

	case 1:
		if m.playersStats == nil {
			sb.WriteString(ui.NormalStyle.Render("Загрузка данных..."))
			break
		}

		// Подвкладки для игроков
		playerTabs := []string{"Золото", "Опыт", "Рейтинг", "Сундуки"}
		for i, tab := range playerTabs {
			if i == m.playersTab {
				sb.WriteString(ui.SelectedStyle.Render(" [" + tab + "] "))
			} else {
				sb.WriteString(ui.NormalStyle.Render(tab + " "))
			}
		}
		sb.WriteString("\n\n")

		sb.WriteString(m.renderPlayersTable())

	case 2:
		if m.guildStats == nil {
			sb.WriteString(ui.NormalStyle.Render("Загрузка данных..."))
			break
		}

		// Подвкладки для гильдий
		guildTabs := []string{"Игроки", "Победы"}
		for i, tab := range guildTabs {
			if i == m.guildsTab {
				sb.WriteString(ui.SelectedStyle.Render(" [" + tab + "] "))
			} else {
				sb.WriteString(ui.NormalStyle.Render(tab + " "))
			}
		}
		sb.WriteString("\n\n")

		sb.WriteString(m.renderGuildsTable())
	}

	if (m.activeTab == 1 && m.playersStats != nil && m.playersStats.PageAmount > 1) ||
		(m.activeTab == 2 && m.guildStats != nil && m.guildStats.PageAmount > 1) {
		var totalPages int
		if m.activeTab == 1 {
			totalPages = m.playersStats.PageAmount
		} else {
			totalPages = m.guildStats.PageAmount
		}
		sb.WriteString(fmt.Sprintf("\nСтраница %d/%d", m.currentPage, totalPages))
	}

	sb.WriteString("\n\n")
	helpText := "←/→ - вкладки"
	if m.activeTab > 0 {
		helpText += ", Tab - подвкладки"
		if m.activeTab == 1 || m.activeTab == 2 {
			helpText += ", ↑/↓ - страницы"
		}
	}
	sb.WriteString(ui.NormalStyle.Render(helpText + ", Esc - назад"))

	return sb.String()
}

func (m *ScoreboardModel) renderPlayersTable() string {
	if m.playersStats == nil || len(m.playersStats.Items) == 0 {
		return "Нет данных для отображения"
	}

	var headers []string
	var widths []int
	var rows [][]string

	// Настройка колонок в зависимости от выбранной подвкладки
	switch m.playersTab {
	case 0: // Золото
		headers = []string{"Позиция", "Игрок", "Золото"}
		widths = []int{10, 30, 15}
		for _, p := range m.playersStats.Items {
			rows = append(rows, []string{
				fmt.Sprintf("%d", p.GoldRatingPos),
				p.Name,
				fmt.Sprintf("%d", p.Gold),
			})
		}
	case 1: // Опыт
		headers = []string{"Позиция", "Игрок", "Опыт"}
		widths = []int{10, 30, 15}
		for _, p := range m.playersStats.Items {
			rows = append(rows, []string{
				fmt.Sprintf("%d", p.ExpRatingPos),
				p.Name,
				fmt.Sprintf("%d", p.Experience),
			})
		}
	case 2: // Рейтинг
		headers = []string{"Позиция", "Игрок", "Рейтинг"}
		widths = []int{10, 30, 15}
		for _, p := range m.playersStats.Items {
			rows = append(rows, []string{
				fmt.Sprintf("%d", p.RatingRatingPos),
				p.Name,
				fmt.Sprintf("%d", p.Rating),
			})
		}
	case 3: // Сундуки
		headers = []string{"Позиция", "Игрок", "Сундуки"}
		widths = []int{10, 30, 15}
		for _, p := range m.playersStats.Items {
			rows = append(rows, []string{
				fmt.Sprintf("%d", p.ChestsOpenedRatingPos),
				p.Name,
				fmt.Sprintf("%d", p.ChestsOpened),
			})
		}
	}

	table := ui.NewTable(m.tableWidth, widths)
	table.AddHeader(headers)
	for _, row := range rows {
		table.AddRow(row)
	}

	return table.Render()
}

func (m *ScoreboardModel) renderGuildsTable() string {
	if m.guildStats == nil || len(m.guildStats.Items) == 0 {
		return "Нет данных для отображения"
	}

	var headers []string
	var widths []int
	var rows [][]string

	// Настройка колонок в зависимости от выбранной подвкладки
	switch m.guildsTab {
	case 0: // Игроки
		headers = []string{"Позиция", "Гильдия", "Игроки"}
		widths = []int{10, 30, 15}
		for _, g := range m.guildStats.Items {
			rows = append(rows, []string{
				fmt.Sprintf("%d", g.GuildMembersRatingPos),
				g.Name,
				fmt.Sprintf("%d", g.GuildMembers),
			})
		}
	case 1: // Победы
		headers = []string{"Позиция", "Гильдия", "Победы"}
		widths = []int{10, 30, 15}
		for _, g := range m.guildStats.Items {
			rows = append(rows, []string{
				fmt.Sprintf("%d", g.WarsVictoriesRatingPos),
				g.Name,
				fmt.Sprintf("%d", g.WarsVictories),
			})
		}
	}

	table := ui.NewTable(m.tableWidth, widths)
	table.AddHeader(headers)
	for _, row := range rows {
		table.AddRow(row)
	}

	return table.Render()
}

func (m *ScoreboardModel) loadStats() tea.Msg {
	ctx := context.Background()

	switch m.activeTab {
	case 0:
		// Получение статистики текущего пользователя
		stats, err := m.Clients.ScoreboardClient.GetCurrentUserStats(ctx, m.id)
		if err != nil {
			return err
		}
		return stats

	case 1:
		// Определение параметра сортировки для игроков
		var orderBy string
		switch m.playersTab {
		case 0:
			orderBy = "gold"
		case 1:
			orderBy = "experience"
		case 2:
			orderBy = "rating"
		case 3:
			orderBy = "chest_opened"
		}

		// Получение списка игроков
		stats, err := m.Clients.ScoreboardClient.GetUserStats(
			ctx,
			nil, // все пользователи
			"",  // без фильтра по имени
			orderBy,
			true, // по убыванию
			pageSize,
			m.currentPage,
		)
		if err != nil {
			return err
		}
		return stats

	case 2:
		// Определение параметра сортировки для гильдий
		var orderBy string
		switch m.guildsTab {
		case 0:
			orderBy = "players"
		case 1:
			orderBy = "wins"
		}

		// Получение списка гильдий
		stats, err := m.Clients.ScoreboardClient.GetGuildStats(
			ctx,
			nil, // все гильдии
			"",  // без фильтра по имени
			orderBy,
			true, // по убыванию
			pageSize,
			m.currentPage,
		)
		if err != nil {
			return err
		}
		return stats
	}

	return nil
}
