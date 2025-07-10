package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"lesta-battleship/cli/internal/api/websocket"
	"lesta-battleship/cli/internal/api/websocket/packets"
	"lesta-battleship/cli/internal/api/websocket/strategies"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/golang-jwt/jwt/v5"
	matchmaking "github.com/lesta-battleship/matchmaking/pkg/packets"
)

var (
	addrFlag = flag.String("addr", "localhost:8080", "")
	pathFlag = flag.String("path", "/matchmaking/custom", "")
)

func main() {
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: *addrFlag, Path: *pathFlag}
	log.Printf("Connecting to %s", u.String())

	id := rand.Text()
	log.Printf("Created ID=%q", id)
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{"sub": id})
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		log.Println(err)

		return
	}
	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	client, err := websocket.NewWebsocketClient(
		u.String(), header, &strategies.MatchmakingStrategy{})
	if err != nil {
		log.Println(err)

		return
	}
	go client.ReadPump()
	go client.WritePump()
	defer client.Stop()
	log.Printf("Connected to %s", u.String())

	done, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		for {
			packet := client.GetPacket()

			var unwrap matchmaking.Packet
			if err := packets.UnwrapAsMatchmaking(packet, &unwrap); err != nil {
				log.Println(err)
			}

			fmt.Println(unwrap.Body)
		}
	}()

	go func() {
		for {
			select {
			default:
				var text string
				fmt.Scanln(&text)

				switch text {
				case "quit":
					packet := matchmaking.NewDisconnect(id)

					client.SendPacket(packets.WrapMatchmaking(packet))

					cancel()
					return
				case "create":
					packet := matchmaking.NewCreateRoom(id)

					client.SendPacket(packets.WrapMatchmaking(packet))
				case "join":
					var roomId string
					fmt.Scanln(&roomId)

					packet := matchmaking.NewJoinRoom(id, roomId)

					client.SendPacket(packets.WrapMatchmaking(packet))
				default:
					packet := matchmaking.NewPlayerMessage(id, text)

					client.SendPacket(packets.WrapMatchmaking(packet))
				}
			case <-done.Done():
				return
			}
		}
	}()

	<-done.Done()
}
