package scoreboard

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lesta-start-battleship/cli/storage/token"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	basePath = ""
)

// Client - клиент для работы с Scoreboard
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	tokenStore *token.Storage
}

// NewClient - создание нового клиента
func NewClient(baseURL string, tokens *token.Storage) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	return &Client{
		baseURL: parsedURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		tokenStore: tokens,
	}, nil
}

// GetUserStats - получение статистики пользователей
func (c *Client) GetUserStats(
	ctx context.Context,
	ids []int,
	nameFilter string,
	orderBy string,
	reverse bool,
	limit int,
	page int,
) (*UserListResponse, error) {
	endpoint := fmt.Sprintf("%susers/", c.baseURL.String())
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("ошибка формирования URL: %w", err)
	}

	// Подготовка параметров запроса
	q := u.Query()
	if ids != nil {
		for _, id := range ids {
			q.Add("ids", strconv.Itoa(id))
		}
	}
	if nameFilter != "" {
		q.Add("name_ilike", nameFilter)
	}
	if reverse {
		q.Add(fmt.Sprintf("order_by_%s", orderBy), "desc")
	}
	q.Add("limit", strconv.Itoa(limit))
	q.Add("page", strconv.Itoa(page))

	u.RawQuery = q.Encode()

	log.Printf("URL запроса рейтинга игроков: %s", u.String())

	body, err := c.doRequest(ctx, u.String())
	if err != nil {
		return nil, err
	}

	var response UserListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return &response, nil
}

// GetCurrentUserStats - получение статистики текущего пользователя
func (c *Client) GetCurrentUserStats(ctx context.Context, username string) (*UserStat, error) {
	resp, err := c.GetUserStats(ctx, nil, username, "", false, 1, 1)
	if err != nil {
		return nil, err
	}

	if len(resp.Items) == 0 {
		return nil, nil
	}

	return &resp.Items[0], nil
}

// GetGuildStats - получение статистики гильдий
func (c *Client) GetGuildStats(
	ctx context.Context,
	guildID *int,
	nameFilter string,
	orderBy string,
	reverse bool,
	limit int,
	page int,
) (*GuildListResponse, error) {
	endpoint := fmt.Sprintf("%sguilds/", c.baseURL.String())
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("ошибка формирования URL: %w", err)
	}

	q := u.Query()
	if guildID != nil {
		q.Add("ids", strconv.Itoa(*guildID))
	}
	if nameFilter != "" {
		q.Add("tag_ilike", nameFilter)
	}
	if reverse {
		q.Add(fmt.Sprintf("order_by_%s", orderBy), "desc")
	}
	q.Add("limit", strconv.Itoa(limit))
	q.Add("page", strconv.Itoa(page))

	u.RawQuery = q.Encode()

	log.Printf("URL запроса рейтинга гильдий: %s", u.String())

	body, err := c.doRequest(ctx, u.String())
	if err != nil {
		return nil, err
	}

	var response GuildListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return &response, nil
}

// doRequest - шаблон для GET запросов и получения ответов
func (c *Client) doRequest(ctx context.Context, url string) ([]byte, error) {
	// Создаем HTTP-запрос
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// установка токена в хедер
	req.Header.Set("Accept", "application/json")
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	// обработка HTTP-ошибок
	if resp.StatusCode >= 400 {
		errorMsg := string(body)

		switch resp.StatusCode {
		case 400:
			return nil, fmt.Errorf("неверный запрос: %s", errorMsg)
		case 401:
			return nil, fmt.Errorf("требуется авторизация: %s", errorMsg)
		case 403:
			return nil, fmt.Errorf("доступ запрещен: %s", errorMsg)
		case 404:
			return nil, fmt.Errorf("ресурс не найден: %s", errorMsg)
		case 500:
			return nil, fmt.Errorf("внутренняя ошибка сервера: %s", errorMsg)
		default:
			return nil, fmt.Errorf("ошибка HTTP %d: %s", resp.StatusCode, errorMsg)
		}
	}

	return body, nil
}
