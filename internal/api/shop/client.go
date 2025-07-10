package shop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"lesta-start-battleship/cli/storage/token"
	"net/http"
	"net/url"
	"time"
)

// Client - клиент для работы с Shop
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
		baseURL:    parsedURL,
		httpClient: &http.Client{Timeout: 15 * time.Second},
		tokenStore: tokens,
	}, nil
}

// doRequest - шаблон для создания запросов
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	reqURL := c.baseURL.ResolveReference(&url.URL{Path: path})

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("error encoding body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), &buf)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	access, refresh := c.tokenStore.GetToken()
	if access != "" {
		req.Header.Set("Authorization", access)
		req.Header.Set("Refresh-Token", refresh)
	}

	return c.httpClient.Do(req)
}

// GetProducts - получение списка предметов
func (c *Client) GetProducts(ctx context.Context) ([]Product, error) {
	resp, err := c.doRequest(ctx, "GET", "item/", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, refresh := c.tokenStore.GetToken()
	if newAccess := resp.Header.Get("Authorization"); newAccess != "" {
		c.tokenStore.SetTokens(newAccess, refresh)
		if newRefresh := resp.Header.Get("Refresh-Token"); newRefresh != "" {
			c.tokenStore.SetTokens(newAccess, newRefresh)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var products []Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return products, nil
}

// GetChests - получение списка сундуков
func (c *Client) GetChests(ctx context.Context) ([]Chest, error) {
	resp, err := c.doRequest(ctx, "GET", "chest/", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, refresh := c.tokenStore.GetToken()
	if newAccess := resp.Header.Get("Authorization"); newAccess != "" {
		c.tokenStore.SetTokens(newAccess, refresh)
		if newRefresh := resp.Header.Get("Refresh-Token"); newRefresh != "" {
			c.tokenStore.SetTokens(newAccess, newRefresh)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var chests []Chest
	if err := json.NewDecoder(resp.Body).Decode(&chests); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return chests, nil
}

// GetPromotions - получение списка акций
func (c *Client) GetPromotions(ctx context.Context) ([]Promotion, error) {
	resp, err := c.doRequest(ctx, "GET", "promotion/", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, refresh := c.tokenStore.GetToken()
	if newAccess := resp.Header.Get("Authorization"); newAccess != "" {
		c.tokenStore.SetTokens(newAccess, refresh)
		if newRefresh := resp.Header.Get("Refresh-Token"); newRefresh != "" {
			c.tokenStore.SetTokens(newAccess, newRefresh)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var promotions []Promotion
	if err := json.NewDecoder(resp.Body).Decode(&promotions); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return promotions, nil
}

// BuyProduct - покупка предмета
func (c *Client) BuyProduct(ctx context.Context, itemID int) error {
	path := fmt.Sprintf("item/%d/buy/", itemID)
	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, refresh := c.tokenStore.GetToken()
	if newAccess := resp.Header.Get("Authorization"); newAccess != "" {
		c.tokenStore.SetTokens(newAccess, refresh)
		if newRefresh := resp.Header.Get("Refresh-Token"); newRefresh != "" {
			c.tokenStore.SetTokens(newAccess, newRefresh)
		}
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

// BuyChest - покупка сундука
func (c *Client) BuyChest(ctx context.Context, chestID int) error {
	path := fmt.Sprintf("chest/%d/buy/", chestID)
	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, refresh := c.tokenStore.GetToken()
	if newAccess := resp.Header.Get("Authorization"); newAccess != "" {
		c.tokenStore.SetTokens(newAccess, refresh)
		if newRefresh := resp.Header.Get("Refresh-Token"); newRefresh != "" {
			c.tokenStore.SetTokens(newAccess, newRefresh)
		}
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

// OpenChest - открытие сундука
func (c *Client) OpenChest(ctx context.Context, chestID, amount int) error {
	requestBody := OpenChestRequest{
		ChestID: chestID,
		Amount:  amount,
	}

	resp, err := c.doRequest(ctx, "POST", "chest/open/", requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, refresh := c.tokenStore.GetToken()
	if newAccess := resp.Header.Get("Authorization"); newAccess != "" {
		c.tokenStore.SetTokens(newAccess, refresh)
		if newRefresh := resp.Header.Get("Refresh-Token"); newRefresh != "" {
			c.tokenStore.SetTokens(newAccess, newRefresh)
		}
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

// TODO:
// 	GetUserPurchases - получить историю покупок пользователя ( GET /purchase/ )
//  Вернет список всех покупок текущего пользователя в разделе "История покупок"

func (c *Client) GetUserPurchases(ctx context.Context) ([]Purchase, error) {
	// Заглушка для будущей реализации
	return nil, fmt.Errorf("not implemented yet")
}

// TODO:
//  GetPromotionDetails - получить детали акции ( GET /promotion/{id}/ )
//  Вернет полную информацию о конкретной акции на её странице

func (c *Client) GetPromotionDetails(ctx context.Context, promotionID int) (*Promotion, error) {
	// Заглушка для будущей реализации
	return nil, fmt.Errorf("not implemented yet")
}
