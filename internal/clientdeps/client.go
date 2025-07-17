package clientdeps

import (
	"lesta-start-battleship/cli/internal/api/auth"
	"lesta-start-battleship/cli/internal/api/guilds"
	"lesta-start-battleship/cli/internal/api/inventory"
	"lesta-start-battleship/cli/internal/api/matchmaking"
	"lesta-start-battleship/cli/internal/api/scoreboard"
	"lesta-start-battleship/cli/internal/api/shop"
)

type Client struct {
	AuthClient       *auth.Client
	GuildsClient     *guilds.Client
	InventoryClient  *inventory.Client
	ScoreboardClient *scoreboard.Client
	ShopClient       *shop.Client
	Matchmaking      *matchmaking.Client
}
