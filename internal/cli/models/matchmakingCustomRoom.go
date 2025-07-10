package models

import (
	"fmt"
	"lesta-start-battleship/cli/internal/api/websocket"
	"lesta-start-battleship/cli/internal/api/websocket/packets"
	"lesta-start-battleship/cli/internal/cli/ui"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	matchmaking "github.com/lesta-battleship/matchmaking/pkg/packets"
)

type MatchmakingCustomRoomModel struct {
	parent   tea.Model
	userId   string
	username string

	wsClient *websocket.WebsocketClient
	roomId   string
}

func NewMatchmakingCustomRoomModel(parent tea.Model, username, userId string, client *websocket.WebsocketClient) *MatchmakingCustomRoomModel {
	return &MatchmakingCustomRoomModel{
		parent:   parent,
		userId:   userId,
		username: username,

		wsClient: client,
		roomId:   "Wait",
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
			packet := matchmaking.NewDisconnect(m.userId)

			m.wsClient.SendPacket(packets.WrapMatchmaking(packet))
			return m.parent, nil

		case tea.KeyCtrlC:
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
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.username))
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
				log.Fatal(err)
			}
			return unwrapped.Body
		}
	}
}
