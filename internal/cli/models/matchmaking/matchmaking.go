package matchmaking

import (
	api "lesta-start-battleship/cli/internal/api/matchmaking"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const matchTypesAmount = 3

const (
	randomSelected int = iota
	rankedSelected
	customSelected
)

type MatchmakingModel struct {
	parent tea.Model

	selected int

	player  *clientdeps.PlayerInfo
	clients *clientdeps.Client
}

func NewMatchmakingModel(parent tea.Model, player *clientdeps.PlayerInfo, clients *clientdeps.Client) *MatchmakingModel {
	return &MatchmakingModel{
		parent: parent,

		player:  player,
		clients: clients,
	}
}

func (m *MatchmakingModel) Init() tea.Cmd {
	m.selected = 0

	return nil
}

func (m *MatchmakingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			m.selected = (m.selected - 1 + matchTypesAmount) % matchTypesAmount
			return m, nil

		case tea.KeyDown:
			m.selected = (m.selected + 1) % matchTypesAmount
			return m, nil

		case tea.KeyEnter:
			switch m.selected {
			case randomSelected:
				model := NewMatchmakingWaitScreenModel(m, api.RandomMatchmaking, m.player, m.clients)

				return model, model.Init()
			case rankedSelected:
				model := NewMatchmakingWaitScreenModel(m, api.RankedMatchmaking, m.player, m.clients)

				return model, model.Init()
			case customSelected:
				model := NewMatchmakingCustomMenuModel(m, m.player, m.clients)

				return model, model.Init()
			}
			return m, nil

		case tea.KeyEsc:
			return m.parent, m.parent.Init()

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	return m, nil
}

var menuItems = [matchTypesAmount]string{
	"Случайный",
	"Рейтинговый",
	"Кастомный",
}

func (m *MatchmakingModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Морской Бой"))
	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.player.Name()))
	sb.WriteString("\n\n")

	for i, item := range menuItems {
		if i == m.selected {
			sb.WriteString(ui.SelectedStyle.Render("> " + item))
		} else {
			sb.WriteString(ui.NormalStyle.Render("  " + item))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(ui.HelpStyle.Render("↑/↓ - выбор, Enter - подтвердить, Esc - выход"))

	return sb.String()
}
