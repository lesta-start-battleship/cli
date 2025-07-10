package strategies

import (
	"fmt"
	"lesta-start-battleship/cli/internal/api/websocket/packets"
	"lesta-start-battleship/cli/internal/api/websocket/packets/guild"

	"github.com/gorilla/websocket"
)

// Стратегия для WebsocketClient.
//
// Ожидает от сервера пакеты типа guild.Packet.
//
// При получении пакета guild.Disconnect принудительно заканчивает работу.
type GuildChatStrategy struct{}

func (c GuildChatStrategy) ReadPump(readChan chan<- packets.Packet, conn *websocket.Conn) error {
	isFirstMessage := true
	for {
		var packet guild.Packet = new(guild.ChatHistoryMessage)
		if isFirstMessage {
			packet = new(guild.ChatHistory)
			isFirstMessage = false
		}

		if err := conn.ReadJSON(&packet); err != nil {
			return fmt.Errorf("GuildChatStrategy.ReadPump: [%w]", err)
		}

		readChan <- packets.WrapGuild(packet)
	}
}

func (c GuildChatStrategy) WritePump(writeChan <-chan packets.Packet, conn *websocket.Conn) error {
	for packet := range writeChan {
		var unwrap guild.Packet
		if err := packets.UnwrapAsGuild(packet, &unwrap); err != nil {
			return fmt.Errorf("GuildChatStrategy.WritePump: [%w]", err)
		}

		if err := conn.WriteJSON(packet.Content()); err != nil {
			return fmt.Errorf("GuildChatStrategy.WritePump: [%w]", err)
		}

		switch unwrap.(type) {
		case *guild.Disconnect:
			return nil
		}
	}

	return nil
}
