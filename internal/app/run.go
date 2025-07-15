package app

import (
	"lesta-start-battleship/cli/internal/api/auth"
	"lesta-start-battleship/cli/internal/api/guilds"
	"lesta-start-battleship/cli/internal/api/inventory"
	"lesta-start-battleship/cli/internal/api/scoreboard"
	"lesta-start-battleship/cli/internal/api/shop"
	cliModel "lesta-start-battleship/cli/internal/cli/initCli"
	"lesta-start-battleship/cli/internal/clientdeps"
	"lesta-start-battleship/cli/storage/token"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	authURL       = "https://battleship-lesta-start.ru/"
	guildsURL     = "https://battleship-lesta-start.ru/guild/"
	guildsWsURL   = "ws://37.9.53.187:8000/api/v1/"
	inventoryURL  = "https://battleship-lesta-start.ru/inventory/"
	scoreboardURL = "https://battleship-lesta-start.ru/scoreboard/"
	shopURL       = "https://battleship-lesta-start.ru/shop/"
)

type App struct {
	program *tea.Program
}

func New() (*App, error) {
	tokenStorage := token.NewStorage()

	initialClients, err := initClients(tokenStorage)
	if err != nil {
		return nil, err
	}

	initialModel := cliModel.NewCLI(initialClients)

	program := tea.NewProgram(initialModel, tea.WithAltScreen())

	return &App{
		program: program,
	}, nil
}

func (a *App) Run() error {
	if _, err := a.program.Run(); err != nil {
		return err
	}
	return nil
}

func initClients(tokenStore *token.Storage) (*clientdeps.Client, error) {
	authClient, err := auth.NewClient(authURL, tokenStore)
	if err != nil {
		return nil, err
	}

	guildsClient, err := guilds.NewClient(guildsURL, guildsWsURL, tokenStore)
	if err != nil {
		return nil, err
	}

	inventoryClient, err := inventory.NewClient(inventoryURL, tokenStore)
	if err != nil {
		return nil, err
	}

	scoreboardClient, err := scoreboard.NewClient(scoreboardURL, tokenStore)
	if err != nil {
		return nil, err
	}

	shopClient, err := shop.NewClient(shopURL, tokenStore)
	if err != nil {
		return nil, err
	}

	return &clientdeps.Client{
		AuthClient:       authClient,
		GuildsClient:     guildsClient,
		InventoryClient:  inventoryClient,
		ScoreboardClient: scoreboardClient,
		ShopClient:       shopClient,
	}, nil
}
