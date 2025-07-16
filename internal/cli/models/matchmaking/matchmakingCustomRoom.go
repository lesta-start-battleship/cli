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

type MatchmakingCustomRoomModel struct {
	parent tea.Model

	roomId string

	player   *clientdeps.PlayerInfo
	wsClient *websocket.WebsocketClient
}

func NewMatchmakingCustomRoomModel(parent tea.Model, player *clientdeps.PlayerInfo, client *websocket.WebsocketClient) *MatchmakingCustomRoomModel {
	return &MatchmakingCustomRoomModel{
		parent: parent,

		roomId: "Wait...",

		player: player,
		wsClient: client,
	}
}

func (m *MatchmakingCustomRoomModel) Init() tea.Cmd {
	return m.waitForMessage()
}

func (m *MatchmakingCustomRoomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			packet := matchmaking.DisconnectPacket(m.player.Id())
			m.wsClient.SendPacket(packets.WrapMatchmaking(packet))

			return m.parent, m.parent.Init()

		case tea.KeyCtrlC:
			packet := matchmaking.DisconnectPacket(m.player.Id())
			m.wsClient.SendPacket(packets.WrapMatchmaking(packet))

			return m, tea.Quit
		}
	case *matchmaking.PlayerMessage:
		m.roomId = msg.Msg

		return m, m.waitForMessage()
	}

	return m, nil
}

func (m *MatchmakingCustomRoomModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Морской Бой"))
	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.player.Name()))
	sb.WriteString("\n\n")

	fmt.Fprintf(&sb, "ID: %q", m.roomId)

	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Esc - выход"))

	return sb.String()
}

func (c *MatchmakingCustomRoomModel) waitForMessage() tea.Cmd {
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
