package matchmaking

import (
	"strconv"

	matchmaking "github.com/lesta-battleship/matchmaking/pkg/packets"
)

type Packet = matchmaking.Packet

type PlayerMessage = matchmaking.PlayerMessage

func CreateRoomPacket(playerId int) Packet {
	senderId := strconv.Itoa(playerId)

	return matchmaking.NewCreateRoom(senderId)
}

func JoinRoomPacket(playerId int, roomId string) Packet {
	senderId := strconv.Itoa(playerId)

	return matchmaking.NewJoinRoom(senderId, roomId)
}

func DisconnectPacket(playerId int) Packet {
	senderId := strconv.Itoa(playerId)

	return matchmaking.NewDisconnect(senderId)
}
