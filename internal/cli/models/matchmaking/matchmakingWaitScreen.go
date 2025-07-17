package matchmaking

import (
	"fmt"
	api "lesta-start-battleship/cli/internal/api/matchmaking"
	"lesta-start-battleship/cli/internal/api/websocket"
	"lesta-start-battleship/cli/internal/api/websocket/packets"
	"lesta-start-battleship/cli/internal/api/websocket/packets/matchmaking"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

type MatchmakingWaitScreenModel struct {
	parent tea.Model

	ticker    *time.Ticker
	startTime time.Time
	endTime   time.Time

	matchPath api.MatchmakingPath
	player    *clientdeps.PlayerInfo
	clients   *clientdeps.Client
	wsClient  *websocket.WebsocketClient
}

func NewMatchmakingWaitScreenModel(parent tea.Model, matchPath api.MatchmakingPath, player *clientdeps.PlayerInfo, clients *clientdeps.Client) *MatchmakingWaitScreenModel {
	now := time.Now()
	ticker := time.NewTicker(time.Second)

	return &MatchmakingWaitScreenModel{
		parent: parent,

		ticker:    ticker,
		startTime: now,
		endTime:   now,

		matchPath: matchPath,
		player:    player,
		clients:   clients,
		wsClient:  nil,
	}
}

func (m *MatchmakingWaitScreenModel) Init() tea.Cmd {
	wsClient := m.wsClient
	if wsClient == nil || wsClient.Connected() {
		client, err := m.clients.Matchmaking.Queue(m.matchPath)
		if err != nil {
		}
		go client.ReadPump()
		go client.WritePump()

		wsClient = client
	}

	now := time.Now()
	ticker := time.NewTicker(time.Second)

	m.ticker = ticker
	m.startTime = now
	m.endTime = now
	m.wsClient = wsClient

	return m.waitForMessage()
}

func (m *MatchmakingWaitScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		model := NewMatchmakingCustomRoomModel(m, m.player, m.wsClient)
		model.roomId = msg.Msg

		return model, model.Init()
	case tickMsg:
		m.endTime = time.Time(msg)
		return m, m.waitForMessage()
	}

	return m, nil
}

func (m *MatchmakingWaitScreenModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Морской Бой"))
	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.player.Name()))
	sb.WriteString("\n\n")

	fmt.Fprintf(&sb, "Время прошло: %s", m.endTime.Sub(m.startTime).Round(time.Second))

	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Esc - выход"))

	return sb.String()
}

func (c *MatchmakingWaitScreenModel) waitForMessage() tea.Cmd {
	return func() tea.Msg {
		select {
		case packet := <-c.wsClient.ReadChan():
			var unwrapped matchmaking.Packet
			if err := packets.UnwrapAsMatchmaking(packet, &unwrapped); err != nil {
				log.Println(err)
			}
			return unwrapped.Body
		case tick := <-c.ticker.C:
			return tickMsg(tick)
		}
	}
}
