package guilds

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lesta-start-battleship/cli/internal/api/websocket"
	"lesta-start-battleship/cli/internal/api/websocket/strategies"
	"lesta-start-battleship/cli/storage/token"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Client - клиент для работы с API гильдий
type Client struct {
	baseURL    *url.URL
	baseWsURL  *url.URL
	httpClient *http.Client
	tokenStore *token.Storage
}

// NewClient создает новый клиент
func NewClient(baseURL string, baseWsURL string, tokens *token.Storage) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	parsedWsURL, err := url.Parse(baseWsURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base WS URL: %w", err)
	}

	return &Client{
		baseURL:   parsedURL,
		baseWsURL: parsedWsURL,

		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		tokenStore: tokens,
	}, nil
}

// doRequest HTTP запрос с заданным методом, путем и телом и с учётом query-параметров
func (c *Client) doRequest(
	ctx context.Context,
	method, path string,
	queryParams map[string]string,
	body interface{},
) ([]byte, error) {
	reqURL := c.baseURL.ResolveReference(&url.URL{Path: path})

	// добавление query-параметров
	if queryParams != nil {
		q := reqURL.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		reqURL.RawQuery = q.Encode()
	}

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("error encoding request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), &buf)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	access, refresh := c.tokenStore.GetToken()
	if access != "" {
		req.Header.Set("Authorization", "Bearer "+access)
		req.Header.Set("Refresh-Token", refresh)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if newAccess := resp.Header.Get("Authorization"); newAccess != "" {
		c.tokenStore.SetTokens(newAccess, refresh)
		if newRefresh := resp.Header.Get("Refresh-Token"); newRefresh != "" {
			c.tokenStore.SetTokens(newAccess, newRefresh)
		}
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

// Устанавливает WebSocket соединение по стандартному URL и переданному Path.
//
// Возвращает ошибку при отсутствие возможности установить соединение.
func (c *Client) doWsConnect(path string) (*websocket.WebsocketClient, error) {
	reqURL := c.baseWsURL.ResolveReference(&url.URL{Path: path})

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	access, refresh := c.tokenStore.GetToken()
	if access != "" {
		header.Set("Authorization", "Bearer "+access)
		header.Set("Refresh-Token", refresh)
	}

	wsClient, err := websocket.NewWebsocketClient(reqURL.String(), header, strategies.GuildChatStrategy{})
	if err != nil {
		return nil, fmt.Errorf("error connecting by websocket: %w", err)
	}

	return wsClient, nil
}

// GetMemberByUserID - получить инфо об участнике по user_id
func (c *Client) GetMemberByUserID(ctx context.Context, userID int) (*MemberResponse, error) {
	path := fmt.Sprintf(PathGetMemberByUserID, userID)
	body, err := c.doRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	var resp ResponseMember
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Value, nil
}

// GetGuildByTag - получить инфо о гильдии по тегу
func (c *Client) GetGuildByTag(ctx context.Context, tag string) (*GuildResponse, error) {
	path := fmt.Sprintf(PathGetGuildByTag, tag)
	body, err := c.doRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	var resp ResponseGuild
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Value, nil
}

// GetGuilds - получить список гильдий с пагинацией
func (c *Client) GetGuilds(ctx context.Context, offset, limit int) (*GuildPagination, error) {
	path := PathGetGuilds
	params := map[string]string{
		"offset": strconv.Itoa(offset),
		"limit":  strconv.Itoa(limit),
	}
	body, err := c.doRequest(ctx, "GET", path, params, nil)
	if err != nil {
		return nil, err
	}
	var resp ResponseGuildPagination
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// SendJoinRequest - отправить запрос на вступление в гильдию
func (c *Client) SendJoinRequest(ctx context.Context, guildTag string, userID int) error {
	path := fmt.Sprintf(PathSendJoinRequest, guildTag)
	params := map[string]string{"user_id": strconv.Itoa(userID)}
	_, err := c.doRequest(ctx, "POST", path, params, nil)
	return err
}

// GetJoinRequests - получить список заявок на вступление (для owner/officer)
func (c *Client) GetJoinRequests(ctx context.Context, guildTag string, userID int) (*RequestPagination, error) {
	path := fmt.Sprintf(PathGetJoinRequests, guildTag)
	params := map[string]string{"user_id": strconv.Itoa(userID)}
	body, err := c.doRequest(ctx, "GET", path, params, nil)
	if err != nil {
		return nil, err
	}
	var resp ResponseRequestPagination
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// ApplyJoinRequest - принять заявку на вступление (owner/officer)
func (c *Client) ApplyJoinRequest(ctx context.Context, guildTag string, userID int, guildMemberID int) error {
	path := fmt.Sprintf(PathApplyJoinRequest, guildTag, userID)
	params := map[string]string{"guild_member_id": strconv.Itoa(guildMemberID)}
	_, err := c.doRequest(ctx, "POST", path, params, nil)
	return err
}

// CancelJoinRequest - отклонить заявку на вступление (owner/officer)
func (c *Client) CancelJoinRequest(ctx context.Context, guildTag string, userID int, guildMemberID int) error {
	path := fmt.Sprintf(PathCancelJoinRequest, guildTag, userID)
	params := map[string]string{"guild_member_id": strconv.Itoa(guildMemberID)}
	_, err := c.doRequest(ctx, "DELETE", path, params, nil)
	return err
}

// CreateGuild - создать новую гильдию
func (c *Client) CreateGuild(ctx context.Context, userID int, req CreateGuildRequest) (*GuildResponse, error) {
	path := PathCreateGuild
	params := map[string]string{"user_id": strconv.Itoa(userID)}
	body, err := c.doRequest(ctx, "POST", path, params, &req)
	if err != nil {
		return nil, err
	}
	var resp ResponseGuild
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// editGuild - изменить данные гильдии (owner)
func (c *Client) EditGuild(ctx context.Context, tag string, userID int, req EditGuildRequest) (*GuildResponse, error) {
	path := fmt.Sprintf(PathEditGuild, tag)
	params := map[string]string{"user_id": strconv.Itoa(userID)}
	body, err := c.doRequest(ctx, "PATCH", path, params, &req)
	if err != nil {
		return nil, err
	}
	var resp ResponseGuild
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// DeleteGuild - удалить свою гильдию (owner)
func (c *Client) DeleteGuild(ctx context.Context, tag string, userID int) error {
	path := fmt.Sprintf(PathDeleteGuild, tag)
	params := map[string]string{"user_id": strconv.Itoa(userID)}
	_, err := c.doRequest(ctx, "DELETE", path, params, nil)
	return err
}

// GetGuildMembers - получить список участников гильдии (с пагинацией)
func (c *Client) GetGuildMembers(ctx context.Context, tag string, offset, limit int) (*MemberPagination, error) {
	path := fmt.Sprintf(PathGetGuildMembers, tag)
	params := map[string]string{
		"offset": strconv.Itoa(offset),
		"limit":  strconv.Itoa(limit),
	}
	body, err := c.doRequest(ctx, "GET", path, params, nil)
	if err != nil {
		return nil, err
	}
	var resp ResponseMemberPagination
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// DeleteMember - удалить участника из гильдии (owner/officer)
func (c *Client) DeleteMember(ctx context.Context, tag string, userID, guildMemberID int) error {
	path := fmt.Sprintf(PathDeleteMember, tag, userID)
	params := map[string]string{"guild_member_id": strconv.Itoa(guildMemberID)}
	_, err := c.doRequest(ctx, "DELETE", path, params, nil)
	return err
}

// EditMember - изменить роль или имя участника (owner/officer)
func (c *Client) EditMember(ctx context.Context, tag string, userID, guildMemberID int, req EditMemberRequest) error {
	path := fmt.Sprintf(PathEditMember, tag, userID)
	params := map[string]string{"guild_member_id": strconv.Itoa(guildMemberID)}
	_, err := c.doRequest(ctx, "PATCH", path, params, &req)
	return err
}

// ExitGuild - выйти из гильдии (любой участник)
func (c *Client) ExitGuild(ctx context.Context, tag string) error {
	path := fmt.Sprintf(PathExitGuild, tag)
	_, err := c.doRequest(ctx, "DELETE", path, nil, nil)
	return err
}

// DeclareWar - объявление войны другой гильдии (вызвать может только владелец гильдии-инициатора)
func (c *Client) DeclareWar(ctx context.Context, initiatorGuildID int,
	targetGuildID int, ownerID int,
) (*DeclareWarResponse, error) {
	reqBody := DeclareWarRequest{
		InitiatorGuildID: initiatorGuildID,
		TargetGuildID:    targetGuildID,
		InitiatorOwnerID: ownerID,
	}

	body, err := c.doRequest(ctx, "POST", PathDeclareWar, nil, reqBody)
	if err != nil {
		return nil, err
	}

	var resp DeclareWarResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ConfirmWar - подтверждение войны гильдией, которую вызывают (вызвать может только владелец гильдии-цели)
func (c *Client) ConfirmWar(ctx context.Context, warID int, targetOwnerID int) (*ConfirmWarResponse, error) {
	path := fmt.Sprintf(PathConfirmWar, warID)
	reqBody := ConfirmWarRequest{TargetOwnerID: targetOwnerID}

	body, err := c.doRequest(ctx, "POST", path, nil, reqBody)
	if err != nil {
		return nil, err
	}

	var resp ConfirmWarResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// CancelWar - отменяет активную или войну в статусе ожидания (вызвать может владелец гильдии-участника с любой стороны)
func (c *Client) CancelWar(ctx context.Context, warID int, ownerID int) (*CancelWarResponse, error) {
	path := fmt.Sprintf(PathCancelWar, warID)
	reqBody := CancelWarRequest{OwnerID: ownerID}

	body, err := c.doRequest(ctx, "POST", path, nil, reqBody)
	if err != nil {
		return nil, err
	}

	var resp CancelWarResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetGuildWarList - получить список войн гильдии с фильтрами
func (c *Client) GetGuildWarList(
	ctx context.Context,
	userID int,
	guildID int,
	isInitiator *bool,
	isTarget *bool,
	status *WarStatus,
	page int,
	pageSize int,
) (*GuildWarListResponse, error) {
	// параметры запроса
	params := map[string]string{
		"user_id":   strconv.Itoa(userID),
		"guild_id":  strconv.Itoa(guildID),
		"page":      strconv.Itoa(page),
		"page_size": strconv.Itoa(pageSize),
	}

	// опциональные параметры
	if isInitiator != nil {
		params["is_initiator"] = strconv.FormatBool(*isInitiator)
	}
	if isTarget != nil {
		params["is_target"] = strconv.FormatBool(*isTarget)
	}
	if status != nil {
		params["status"] = string(*status)
	}

	body, err := c.doRequest(ctx, "GET", PathListGuildWars, params, nil)
	if err != nil {
		return nil, err
	}

	var response GuildWarListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Устанавливает WebSocket соединение с Гильдейским чатом по переданному GuildId.
//
// Возвращает ошибку при отсутствие возможности установить соединение.
func (c *Client) JoinGuildChat(guildId int) (*websocket.WebsocketClient, error) {
	path := fmt.Sprintf(PathConnectGuildChat, guildId)

	client, err := c.doWsConnect(path)
	if err != nil {
		return nil, err
	}

	return client, nil
}
