package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader   = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	serverOnce sync.Once
	chatServer *ChatServer
)

type ChatServer struct {
	rooms   map[int]*ChatRoom
	roomsMu sync.Mutex
}

type ChatRoom struct {
	clients  map[*websocket.Conn]string
	messages []ChatMessage
	mu       sync.Mutex
}

type ChatMessage struct {
	GuildID   int       `json:"guild_id"`
	Username  string    `json:"username"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

func NewChatServer() *ChatServer {
	serverOnce.Do(func() {
		chatServer = &ChatServer{rooms: make(map[int]*ChatRoom)}
	})
	return chatServer
}

func (s *ChatServer) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	_, msgBytes, err := conn.ReadMessage()
	if err != nil {
		log.Println("Read error:", err)
		return
	}

	var init struct {
		GuildID  int    `json:"guild_id"`
		Username string `json:"username"`
	}
	if err := json.Unmarshal(msgBytes, &init); err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	log.Printf("User %s connected to guild %d", init.Username, init.GuildID)

	// Получаем или создаем комнату
	s.roomsMu.Lock()
	room, exists := s.rooms[init.GuildID]
	if !exists {
		room = &ChatRoom{
			clients:  make(map[*websocket.Conn]string),
			messages: make([]ChatMessage, 0, 100), // увеличил размер слайса
		}
		s.rooms[init.GuildID] = room
	}
	s.roomsMu.Unlock()

	// Добавляем клиента
	room.mu.Lock()
	room.clients[conn] = init.Username
	room.mu.Unlock()

	// Отправляем историю
	room.mu.Lock()
	for _, msg := range room.messages {
		if err := conn.WriteJSON(msg); err != nil {
			log.Println("Send history error:", err)
		}
	}
	room.mu.Unlock()

	// Читаем сообщения
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg ChatMessage
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Println("Message parse error:", err)
			continue
		}

		log.Printf("Message from %s in guild %d: %s", msg.Username, init.GuildID, msg.Text)

		room.mu.Lock()
		room.messages = append(room.messages, msg)
		if len(room.messages) > 10 {
			room.messages = room.messages[1:]
		}

		for client := range room.clients {
			if err := client.WriteJSON(msg); err != nil {
				delete(room.clients, client)
				client.Close()
			}
		}
		room.mu.Unlock()
	}

	// Удаляем клиента
	room.mu.Lock()
	delete(room.clients, conn)
	room.mu.Unlock()
	log.Printf("User %s disconnected from guild %d", init.Username, init.GuildID)
}

func (r *ChatRoom) messageExists(msg ChatMessage) bool {
	for _, m := range r.messages {
		if m.Username == msg.Username && m.Text == msg.Text && m.Timestamp.Equal(msg.Timestamp) {
			return true
		}
	}
	return false
}

func StartChatServer() {
	server := NewChatServer()
	http.HandleFunc("/ws", server.HandleConnection)
	log.Println("Chat server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	StartChatServer()
}
