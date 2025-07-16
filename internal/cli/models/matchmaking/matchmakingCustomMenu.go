package matchmaking

import (
	api "lesta-start-battleship/cli/internal/api/matchmaking"
	"lesta-start-battleship/cli/internal/api/websocket"
	"lesta-start-battleship/cli/internal/api/websocket/packets"
	"lesta-start-battleship/cli/internal/api/websocket/packets/matchmaking"
	"lesta-start-battleship/cli/internal/cli/handlers"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	createSelected int = iota
	joinSelected
)

type MatchmakingCustomMenuModel struct {
	parent tea.Model

	selected int

	player   *clientdeps.PlayerInfo
	clients  *clientdeps.Client
	wsClient *websocket.WebsocketClient
}

func NewMatchmakingCustomMenuModel(parent tea.Model, player *clientdeps.PlayerInfo, clients *clientdeps.Client) *MatchmakingCustomMenuModel {
	return &MatchmakingCustomMenuModel{
		parent: parent,

		player:   player,
		clients:  clients,
		wsClient: nil,
	}
}

func (m *MatchmakingCustomMenuModel) Init() tea.Cmd {
	wsClient := m.wsClient
	if wsClient == nil || !wsClient.Connected() {
		client, err := m.clients.Matchmaking.Queue(api.CustomMatchmaking)
		if err != nil {
			return func() tea.Msg {
				return handlers.WsErrorMsg{Err: err}
			}
		}
		go client.ReadPump()
		go client.WritePump()

		wsClient = client
	}

	m.selected = 0
	m.wsClient = wsClient

	return nil
}

func (m *MatchmakingCustomMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			case createSelected:
				packet := matchmaking.CreateRoomPacket(m.player.Id())
				m.wsClient.SendPacket(packets.WrapMatchmaking(packet))

				child := NewMatchmakingCustomRoomModel(m, m.player, m.wsClient)

				return child, child.Init()
			case joinSelected:
				child := NewMatchmakingCustomJoinModel(m, m.player, m.wsClient)

				return child, child.Init()
			}
			return m, nil

		case tea.KeyEsc:
			packet := matchmaking.DisconnectPacket(m.player.Id())
			m.wsClient.SendPacket(packets.WrapMatchmaking(packet))

			return m.parent, m.parent.Init()

		case tea.KeyCtrlC:
			packet := matchmaking.DisconnectPacket(m.player.Id())
			m.wsClient.SendPacket(packets.WrapMatchmaking(packet))

			return m, tea.Quit
		}
	case handlers.WsErrorMsg:
		return m.parent, m.parent.Init()
	}

	return m, nil
}

func (m *MatchmakingCustomMenuModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Морской Бой"))
	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.player.Name()))
	sb.WriteString("\n\n")

	menuItems := []string{
		"Создать",
		"Присоединиться",
	}

	for i, item := range menuItems {
		if i == m.selected {
			sb.WriteString(ui.SelectedStyle.Render("> " + item))
		} else {
			sb.WriteString(ui.NormalStyle.Render("  " + item))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(ui.NormalStyle.Render("↑/↓ - выбор, Enter - подтвердить, Esc - выход"))

	return sb.String()
}
