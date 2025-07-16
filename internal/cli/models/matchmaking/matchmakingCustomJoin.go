package matchmaking

import (
	"fmt"
	"lesta-start-battleship/cli/internal/api/websocket"
	"lesta-start-battleship/cli/internal/api/websocket/packets"
	"lesta-start-battleship/cli/internal/api/websocket/packets/matchmaking"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type MatchmakingCustomJoinModel struct {
	parent tea.Model

	input string

	player   *clientdeps.PlayerInfo
	wsClient *websocket.WebsocketClient
}

func NewMatchmakingCustomJoinModel(parent tea.Model, player *clientdeps.PlayerInfo, wsClient *websocket.WebsocketClient) *MatchmakingCustomJoinModel {
	return &MatchmakingCustomJoinModel{
		parent: parent,

		player:   player,
		wsClient: wsClient,
	}
}

func (m *MatchmakingCustomJoinModel) Init() tea.Cmd {
	m.input = ""

	return m.waitForMessage()
}

func (m *MatchmakingCustomJoinModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			m.input += string(msg.Runes)

			return m, m.waitForMessage()
		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case tea.KeyEnter:
			if m.input == "" {
				return m, m.waitForMessage()
			}
			packet := matchmaking.JoinRoomPacket(m.player.Id(), m.input)
			m.wsClient.SendPacket(packets.WrapMatchmaking(packet))

			m.input = ""

			return m, m.waitForMessage()

		case tea.KeyEsc:
			return m.parent, m.parent.Init()
		case tea.KeyCtrlC:
			packet := matchmaking.DisconnectPacket(m.player.Id())
			m.wsClient.SendPacket(packets.WrapMatchmaking(packet))

			return m, tea.Quit
		}

	case *matchmaking.PlayerMessage:
		model := NewMatchmakingCustomRoomModel(m.parent, m.player, m.wsClient)
		model.roomId = msg.Msg

		return model, model.Init()
	}

	return m, nil
}

func (m *MatchmakingCustomJoinModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Морской Бой"))
	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.player.Name()))
	sb.WriteString("\n\n")

	fmt.Fprintf(&sb, "Введите ID: %q", m.input)

	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("↑/↓ - выбор, Enter - подтвердить, Esc - выход"))

	return sb.String()
}

func (c *MatchmakingCustomJoinModel) waitForMessage() tea.Cmd {
	return func() tea.Msg {
		select {
		case packet := <-c.wsClient.ReadChan():
			var unwrapped matchmaking.Packet
			if err := packets.UnwrapAsMatchmaking(packet, &unwrapped); err != nil {
				log.Println(err)
			}
			return unwrapped.Body
		}
	}
}
