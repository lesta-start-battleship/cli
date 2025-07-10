package handlers

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WsConnectedMsg struct{}
type WsDisconnectedMsg struct{}
type ReconnectMsg struct{}
type PingMsg struct{}
type WsErrorMsg struct {
	Err error
}
type CloseChatMsg struct{}

type ChatMessage struct {
	GuildID   int       `json:"guild_id"`
	Username  string    `json:"username"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	IsOwn     bool      `json:"-"`
	IsSystem  bool      `json:"-"`
}

type WsClient struct {
	conn      *websocket.Conn
	Incoming  chan ChatMessage
	Outgoing  chan ChatMessage
	closeChan chan struct{}
	connected bool
	mu        sync.Mutex
	closeOnce sync.Once
}

func NewWsClient() *WsClient {
	return &WsClient{
		Incoming:  make(chan ChatMessage, 100),
		Outgoing:  make(chan ChatMessage, 10),
		closeChan: make(chan struct{}),
	}
}

func (c *WsClient) Connect(guildID int, username string) error {
	/*c.reconnectMu.Lock()
	defer c.reconnectMu.Unlock()*/

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return nil
	}

	c.closeChan = make(chan struct{})

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		return err
	}

	initMsg := struct {
		GuildID  int    `json:"guild_id"`
		Username string `json:"username"`
	}{guildID, username}

	if err := conn.WriteJSON(initMsg); err != nil {
		conn.Close()
		return err
	}

	c.conn = conn
	c.connected = true

	// Запускаем обработчики
	go c.readMessages()
	go c.writeMessages()

	return nil
}

func (c *WsClient) readMessages() {
	defer func() {
		c.mu.Lock()
		if c.conn != nil {
			c.conn.Close()
		}
		c.connected = false
		close(c.Incoming)
		c.mu.Unlock()
	}()

	for {
		select {
		case <-c.closeChan:
			return
		default:
			_, msgBytes, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}

			var msg ChatMessage
			if err := json.Unmarshal(msgBytes, &msg); err != nil {
				log.Println("Failed to parse message:", err)
				continue
			}

			/*select {
			case c.Incoming <- msg:
			case <-c.closeChan:
				return
			}*/
			c.Incoming <- msg
		}
	}
}

func (c *WsClient) writeMessages() {
	defer func() {
		c.mu.Lock()
		if c.conn != nil {
			c.conn.Close()
		}
		c.connected = false
		c.mu.Unlock()
	}()

	for {
		select {
		case <-c.closeChan:
			return
		case msg := <-c.Outgoing:
			c.mu.Lock()
			err := c.conn.WriteJSON(msg)
			c.mu.Unlock()

			if err != nil {
				log.Println("Failed to send message:", err)
				return
			}
		}
	}
}

func (c *WsClient) Close() {
	c.closeOnce.Do(func() {
		close(c.closeChan)

		c.mu.Lock()
		defer c.mu.Unlock()

		if c.conn != nil {
			c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			c.conn.Close()
			c.conn = nil
		}
		c.connected = false
	})
}

func (c *WsClient) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connected
}
