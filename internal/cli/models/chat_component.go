package models

import (
	"fmt"
	"lesta-start-battleship/cli/internal/api/websocket"
	"lesta-start-battleship/cli/internal/api/websocket/packets"
	"lesta-start-battleship/cli/internal/api/websocket/packets/guild"
	"lesta-start-battleship/cli/internal/cli/handlers"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ChatComponent struct {
	Username     string
	guildID      int
	messages     []*guild.ChatHistoryMessage
	input        string
	Focused      bool
	scrollOffset int
	Visible      bool
	Width        int
	err          error
	clients      *clientdeps.Client
	wsClient     *websocket.WebsocketClient
}

func NewChatComponent(username string, guildID int, clients *clientdeps.Client) *ChatComponent {
	return &ChatComponent{
		Username: username,
		guildID:  guildID,
		Width:    55,

		clients: clients,
	}
}

func (c *ChatComponent) Init() tea.Cmd {
	if !c.Visible {
		return nil
	}

	client, err := c.clients.GuildsClient.JoinGuildChat(c.guildID)
	if err != nil {
		return func() tea.Msg {
			return handlers.WsErrorMsg{Err: err}
		}
	}
	go client.ReadPump()
	go client.WritePump()

	c.wsClient = client

	return c.waitForMessage()
}

func (c *ChatComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.Width = msg.Width / 4
		return c, nil

	case *guild.ChatHistory:
		for _, msg := range msg.Data {
			c.messages = append(c.messages, &msg)
		}
		c.scrollToBottom()
		return c, c.waitForMessage()
	case *guild.ChatHistoryMessage:
		c.messages = append(c.messages, msg)
		c.scrollToBottom()
		return c, c.waitForMessage()

	case handlers.WsConnectedMsg:
		return c, c.waitForMessage()

	case handlers.WsErrorMsg:
		c.err = msg.Err
		return c, tea.Tick(5*time.Second, func(time.Time) tea.Msg {
			return handlers.ReconnectMsg{}
		})

	case handlers.ReconnectMsg:
		client, err := c.clients.GuildsClient.JoinGuildChat(c.guildID)
		if err != nil {
			c.err = err
		} else {
			go client.ReadPump()
			go client.WritePump()

			c.wsClient = client
		}

		return c, c.waitForMessage()

	case tea.KeyMsg:
		if !c.Focused {
			return c, nil
		}

		switch msg.Type {
		case tea.KeyEnter:
			if c.input == "" {
				return c, nil
			}
			newMsg := packets.WrapGuild(guild.ChatMessage{Msg: c.input})
			c.input = ""
			c.scrollToBottom()
			c.wsClient.WriteChan() <- newMsg
			return c, tea.Batch(c.waitForMessage(),
				func() tea.Msg { return ChatKeyHandledMsg{} },
			)

		case tea.KeyBackspace:
			if len(c.input) > 0 {
				c.input = c.input[:len(c.input)-1]
			}

		case tea.KeySpace:
			c.input += " "

		case tea.KeyRunes:
			c.input += string(msg.Runes)

		case tea.KeyDown:
			if c.scrollOffset < len(c.messages)-10 {
				c.scrollOffset++
			}
			return c, func() tea.Msg { return ChatKeyHandledMsg{} }

		case tea.KeyUp:
			if c.scrollOffset > 0 {
				c.scrollOffset--
			}
			return c, func() tea.Msg { return ChatKeyHandledMsg{} }

		case tea.KeyEsc:
			c.Close()
			return c, tea.Sequence(
				func() tea.Msg { return ChatClosedMsg{} },
				func() tea.Msg { return ChatKeyHandledMsg{} },
			)
		}
	}

	return c, cmd
}

func (c *ChatComponent) View() string {
	if !c.Visible {
		return ""
	}

	var sb strings.Builder

	header := ui.ChatHeaderStyle.Render(fmt.Sprintf(" Чат гильдии "))
	if c.Focused {
		header = ui.SelectedStyle.Render(fmt.Sprintf(" Чат гильдии (активен) "))
	}

	sb.WriteString(header)
	sb.WriteString("\n\n")

	start := c.scrollOffset
	if start < 0 {
		start = 0
	}
	end := start + 10
	if end > len(c.messages) {
		end = len(c.messages)
	}

	for _, msg := range c.messages[start:end] {
		if msg.Username == c.Username {
			sb.WriteString(ui.OwnMessageStyle.Render(fmt.Sprintf("%s: %s", msg.Username, msg.Content)))
		} else if c.Focused {
			sb.WriteString(ui.OtherMessageStyle.Render(fmt.Sprintf("%s: %s", msg.Username, msg.Content)))
		} else {
			sb.WriteString(ui.NewMessageStyle.Render(fmt.Sprintf("%s: %s", msg.Username, msg.Content)))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	if c.Focused {
		sb.WriteString(ui.ChatInputStyle.Render("> " + c.input + "_"))
	} else {
		sb.WriteString(ui.HelpStyle.Render("Нажмите Ctrl+G для ввода"))
	}

	if c.err != nil {
		sb.WriteString("\n\n")
		sb.WriteString(ui.ErrorStyle.Render("Ошибка: " + c.err.Error()))
	}

	return ui.ChatContainerStyle.Width(c.Width).Render(sb.String())
}

func (c *ChatComponent) waitForMessage() tea.Cmd {
	return func() tea.Msg {
		select {
		case packet := <-c.wsClient.ReadChan():
			var unwrapped guild.Packet
			if err := packets.UnwrapAsGuild(packet, &unwrapped); err != nil {
				log.Println(err)
			}
			return unwrapped
		case <-time.After(30 * time.Second):
			return handlers.PingMsg{}
		}
	}
}

// func (c *ChatComponent) messageExists(msg guild.ChatHistoryMessage) bool {
// 	for _, m := range c.messages {
// 		if m.Username == msg.Username && m.Text == msg.Text && m.Timestamp.Equal(msg.Timestamp) {
// 			return true
// 		}
// 	}
// 	return false
// }

func (c *ChatComponent) scrollToBottom() {
	c.scrollOffset = 0
	if len(c.messages) > 10 {
		c.scrollOffset = len(c.messages) - 10
	}
}

func (c *ChatComponent) Close() {
	if c.wsClient != nil {
		c.wsClient.Stop()
	}
	c.Visible = false
	c.Focused = false
}

func (c *ChatComponent) Toggle() {
	c.Visible = !c.Visible
	if c.Visible {
		c.Focused = true
		if len(c.messages) == 0 {
			c.Init()
		}
	} else {
		c.Focused = false
	}
}

func (c *ChatComponent) Focus() {
	c.Focused = true
}

func (c *ChatComponent) Blur() {
	c.Focused = false
}

func (c *ChatComponent) IsVisible() bool {
	return c.Visible
}
