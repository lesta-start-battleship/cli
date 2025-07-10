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

	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang-jwt/jwt/v5"
	matchmaking "github.com/lesta-battleship/matchmaking/pkg/packets"
)

type MatchmakingCustomMenuModel struct {
	parent   tea.Model
	userId   string
	username string
	selected int

	wsClient *websocket.WebsocketClient
}

func NewMatchmakingCustomMenuModel(parent tea.Model, username string) *MatchmakingCustomMenuModel {
	id := rand.Text()
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{"sub": id})
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		log.Fatal(err)
	}
	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	url := formatMatchmakingUrl("custom")
	client, err := websocket.NewWebsocketClient(url, header, strategies.MatchmakingStrategy{})
	if err != nil {
		log.Fatal(err)
	}
	go client.WritePump()
	go client.ReadPump()

	return &MatchmakingCustomMenuModel{
		parent:   parent,
		userId:   id,
		username: username,
		selected: 0,

		wsClient: client,
	}
}

func (m *MatchmakingCustomMenuModel) Init() tea.Cmd {
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
			case 0:
				packet := matchmaking.NewCreateRoom(m.userId)
				m.wsClient.SendPacket(packets.WrapMatchmaking(packet))

				model := NewMatchmakingCustomRoomModel(m, m.username, m.userId, m.wsClient)
				return model, model.Init()
			case 1:
				model := NewMatchmakingCustomJoinModel(m, m.username, m.userId, m.wsClient)
				return model, model.Init()
			}
			return m, nil

		case tea.KeyEsc:
			packet := matchmaking.NewDisconnect(m.userId)

			m.wsClient.SendPacket(packets.WrapMatchmaking(packet))
			return m.parent, nil

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *MatchmakingCustomMenuModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Морской Бой"))
	sb.WriteString("\n\n")
	sb.WriteString(ui.NormalStyle.Render("Пользователь: " + m.username))
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

func (c *MatchmakingCustomMenuModel) waitForMessage() tea.Cmd {
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
