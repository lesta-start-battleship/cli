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

type MatchmakingCustomJoinModel struct {
	parent   tea.Model
	userId   string
	username string

	input    string
	wsClient *websocket.WebsocketClient
}

func NewMatchmakingCustomJoinModel(parent tea.Model, username, userId string, wsClient *websocket.WebsocketClient) *MatchmakingCustomJoinModel {
	return &MatchmakingCustomJoinModel{
		parent:   parent,
		userId:   userId,
		username: username,

		input:    "",
		wsClient: wsClient,
	}
}

func (m *MatchmakingCustomJoinModel) Init() tea.Cmd {
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
				return m, nil
			}
			newMsg := packets.WrapMatchmaking(matchmaking.NewJoinRoom(m.userId, m.input))
			m.wsClient.WriteChan() <- newMsg
			m.input = ""

			return m, m.waitForMessage()

		case tea.KeyEsc:
			return m.parent, nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case *matchmaking.PlayerMessage:
		model := NewMatchmakingCustomRoomModel(m.parent, m.username, m.userId, m.wsClient)
		model.roomId = msg.Msg
		return model, model.Init()
	}

	return m, nil
}

func (m *MatchmakingCustomJoinModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Морской Бой"))
	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.username))
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
				log.Fatal(err)
			}
			return unwrapped.Body
		}
	}
}
