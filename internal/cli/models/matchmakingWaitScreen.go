package models

import (
	"crypto/rand"
	"fmt"
	"lesta-start-battleship/cli/internal/api/websocket"
	"lesta-start-battleship/cli/internal/api/websocket/packets"
	"lesta-start-battleship/cli/internal/api/websocket/strategies"
	"lesta-start-battleship/cli/internal/cli/ui"
	"log"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang-jwt/jwt/v5"
	matchmaking "github.com/lesta-battleship/matchmaking/pkg/packets"
)

type tickMsg time.Time

type MatchmakingWaitScreenModel struct {
	parent   tea.Model
	userId   string
	username string

	ticker    *time.Ticker
	startTime time.Time
	endTime   time.Time

	wsClient *websocket.WebsocketClient
}

func NewMatchmakingWaitScreenModel(parent tea.Model, username, matchType string) *MatchmakingWaitScreenModel {
	id := rand.Text()
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{"sub": id})
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		log.Fatal(err)
	}
	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	url := formatMatchmakingUrl(matchType)
	client, err := websocket.NewWebsocketClient(url, header, strategies.MatchmakingStrategy{})
	if err != nil {
		log.Fatal(err)
	}
	go client.WritePump()
	go client.ReadPump()

	now := time.Now()
	ticker := time.NewTicker(time.Second)

	return &MatchmakingWaitScreenModel{
		parent:   parent,
		userId:   id,
		username: username,

		ticker:    ticker,
		startTime: now,
		endTime:   now,

		wsClient: client,
	}
}

func (m *MatchmakingWaitScreenModel) Init() tea.Cmd {
	return m.waitForMessage()
}

func (m *MatchmakingWaitScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m.parent, nil

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case *matchmaking.PlayerMessage:
		model := NewMatchmakingCustomRoomModel(m, m.username, m.userId, m.wsClient)
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
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.username))
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
				log.Fatal(err)
			}
			return unwrapped.Body
		case tick := <-c.ticker.C:
			return tickMsg(tick)
		}
	}
}
