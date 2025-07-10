package inventory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lesta-start-battleship/cli/storage/token"
	"net/http"
	"net/url"
	"time"
)

// Client - клиент для взаимодействия с API инвентаря
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	tokenStore *token.Storage
}

// NewClient создает новый клиент для работы с API инвентаря
func NewClient(baseURL string, tokens *token.Storage) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("некорректный базовый URL: %w", err)
	}

	return &Client{
		baseURL: parsedURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second, // Таймаут для безопасности
		},
		tokenStore: tokens,
	}, nil
}

// doRequest выполняет HTTP запрос
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	reqURL := c.baseURL.ResolveReference(&url.URL{Path: path})

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("ошибка кодирования тела запроса: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), &buf)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	access, refresh := c.tokenStore.GetToken()
	if access != "" {
		req.Header.Set("Authorization", access)
		req.Header.Set("Refresh-Token", refresh)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if newAccess := resp.Header.Get("Authorization"); newAccess != "" {
		c.tokenStore.SetTokens(newAccess, refresh)
		if newRefresh := resp.Header.Get("Refresh-Token"); newRefresh != "" {
			c.tokenStore.SetTokens(newAccess, newRefresh)
		}
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		var errResp ErrorResponse
		if json.Unmarshal(body, &errResp) == nil && errResp.Error != "" {
			return nil, fmt.Errorf("ошибка API: %s", errResp.Error)
		}
		return nil, fmt.Errorf("HTTP ошибка %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// GetUserInventory получает инвентарь пользователя
func (c *Client) GetUserInventory(ctx context.Context) (*UserInventoryResponse, error) {
	body, err := c.doRequest(ctx, "GET", UserInventoryPath, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения инвентаря: %w", err)
	}

	var resp UserInventoryResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}
	return &resp, nil
}
