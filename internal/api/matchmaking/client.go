package matchmaking

import (
	"fmt"
	"lesta-start-battleship/cli/internal/api/websocket"
	"lesta-start-battleship/cli/internal/api/websocket/strategies"
	"lesta-start-battleship/cli/storage/token"
	"net/http"
)

// Client - клиент для работы с API матчмейкинга
type Client struct {
	tokenStore *token.Storage
}

// NewClient создает новый клиент
func NewClient(tokens *token.Storage) *Client {
	return &Client{
		tokenStore: tokens,
	}
}

// Устанавливает WebSocket соединение с Матчмейкингом по переданному MatchPath.
//
// Возвращает ошибку при отсутствие возможности установить соединение.
func (c *Client) Queue(matchPath MatchmakingPath) (*websocket.WebsocketClient, error) {
	client, err := c.doWsConnect(matchPath)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Устанавливает WebSocket соединение по стандартному URL и переданному Path.
//
// Возвращает ошибку при отсутствие возможности установить соединение.
func (c *Client) doWsConnect(path MatchmakingPath) (*websocket.WebsocketClient, error) {
	reqURL := formatMatchmakingPath(path)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	access, refresh := c.tokenStore.GetToken()
	if access != "" {
		header.Set("Authorization", "Bearer "+access)
		header.Set("Refresh-Token", refresh)
	}

	wsClient, err := websocket.NewWebsocketClient(reqURL, header, strategies.MatchmakingStrategy{})
	if err != nil {
		return nil, fmt.Errorf("error connecting by websocket: %w", err)
	}

	return wsClient, nil
}
