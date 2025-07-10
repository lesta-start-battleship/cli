package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"lesta-battleship/cli/internal/api/websocket"
	"lesta-battleship/cli/internal/api/websocket/packets"
	"lesta-battleship/cli/internal/api/websocket/packets/guild"
	"lesta-battleship/cli/internal/api/websocket/strategies"
	"log"
	"os"
	"os/signal"
)

var pathFlag = flag.String("url", "ws://localhost:8080", "URL to connect to")

func main() {
	flag.Parse()

	u := *pathFlag
	log.Printf("Connecting to %s", u)
	client, err := websocket.NewWebsocketClient(u, nil, strategies.GuildChatStrategy{})
	if err != nil {
		return
	}
	go client.ReadPump()
	go client.WritePump()
	log.Printf("Connected to %s", u)

	done, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		for {
			packet := client.GetPacket()

			var unwrap guild.Packet
			packets.UnwrapAsGuild(packet, &unwrap)

			message, _ := json.Marshal(unwrap)

			fmt.Printf("%v\n", string(message))
		}
	}()

	go func() {
		for {
			var text string
			fmt.Scanln(&text)
			switch text {
			case "quit":
				chatMsg := packets.WrapGuild(&guild.Disconnect{})

				client.SendPacket(chatMsg)
			default:
				chatMsg := packets.WrapGuild(&guild.ChatMessage{Msg: text})

				client.SendPacket(chatMsg)
			}
		}
	}()

	<-done.Done()
	client.Stop()
}
