package strategies

import (
	"fmt"
	"lesta-start-battleship/cli/internal/api/websocket/packets"

	"github.com/gorilla/websocket"
	matchmaking "github.com/lesta-battleship/matchmaking/pkg/packets"
)

// Стратегия для WebsocketClient.
//
// Ожидает от сервера пакеты типа matchmaking.Packet.
//
// При получении пакета matchmaking.Disconnect заканчивает работу.
type MatchmakingStrategy struct{}

func (c MatchmakingStrategy) ReadPump(readChan chan<- packets.Packet, conn *websocket.Conn) error {
	for {
		packet := matchmaking.Packet{}

		if err := conn.ReadJSON(&packet); err != nil {
			return fmt.Errorf("MatchmakingStrategy.ReadPump: [%w]", err)
		}

		readChan <- packets.WrapMatchmaking(packet)
	}
}

func (c MatchmakingStrategy) WritePump(writeChan <-chan packets.Packet, conn *websocket.Conn) error {
	for packet := range writeChan {
		unwrap := matchmaking.Packet{}
		if err := packets.UnwrapAsMatchmaking(packet, &unwrap); err != nil {
			return fmt.Errorf("MatchmakingStrategy.WritePump: [%w]", err)
		}

		if err := conn.WriteJSON(unwrap); err != nil {
			return fmt.Errorf("MatchmakingStrategy.WritePump: [%w]", err)
		}

		switch unwrap.Body.(type) {
		case *matchmaking.Disconnect:
			return nil
		}
	}

	return nil
}
