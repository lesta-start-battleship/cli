package guilds

import "time"

// константы путей API
const (
	PathGetMemberByUserID = "member/member/%d"
	PathGetGuildByTag     = "%s"
	PathSendJoinRequest   = "request/%s"
	PathGetJoinRequests   = "request/%s"
	PathApplyJoinRequest  = "request/%s/%d/apply"
	PathCancelJoinRequest = "request/%s/%d/cancel"
	PathCreateGuild       = ""
	PathDeleteGuild       = "%s"
	PathGetGuildMembers   = "member/%s"
	PathDeleteMember      = "member/%s/%d"
	PathEditMember        = "member/%s/%d"
	PathExitGuild         = "member/%s"
	PathDeclareWar        = "war/declare"
	PathConfirmWar        = "war/confirm/%d"
	PathCancelWar         = "war/cancel/%d"
	PathListGuildWars     = "war/list"
	PathGetGuilds         = ""
	PathEditGuild         = "%s"
	PathConnectGuildChat  = "chat/ws/guild/%d"
)

// Role - роль участника гильдии
type Role struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	// RolePromote []int  `json:"role_promote"` // список ID ролей, которыми может управлять пользователь
	Permissions []string `json:"permissions"`
}

// MemberResponse - информация об участнике гильдии
type MemberResponse struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	GuildID  int    `json:"guild_id"`
	GuildTag string `json:"guild_tag"`
	Role     Role   `json:"role"`
}

// GuildResponse - информация о гильдии
type GuildResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
	OwnerID     int    `json:"owner_id"`
	IsActive    bool   `json:"is_active"` // активна ли
	IsFull      bool   `json:"is_full"`   // заполнена ли
}

// GuildPagination - список гильдий с пагинацией
type GuildPagination struct {
	Items      []GuildResponse `json:"items"`
	TotalItems int             `json:"total_items"`
	TotalPages int             `json:"total_pages"`
}

// MemberPagination - список участников с пагинацией
type MemberPagination struct {
	Items      []MemberResponse `json:"items"`
	TotalItems int              `json:"total_items"`
	TotalPages int              `json:"total_pages"`
}

// BaseResponse - базовый формат ответа
type BaseResponse struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"error_code"`
	Status    bool   `json:"status"`
}

// ResponseGuild - ответ с данными гильдии
type ResponseGuild struct {
	BaseResponse
	Value *GuildResponse `json:"value"`
}

// ResponseGuildPagination - ответ со списком гильдий
type ResponseGuildPagination struct {
	BaseResponse
	Value *GuildPagination `json:"value"`
}

// ResponseMember - ответ с данными участника
type ResponseMember struct {
	BaseResponse
	Value *MemberResponse `json:"value"`
}

// ResponseMemberPagination - ответ со списком участников
type ResponseMemberPagination struct {
	BaseResponse
	Value *MemberPagination `json:"value"`
}

// RequestResponse - заявка на вступление
type RequestResponse struct {
	UserID    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	CreatedAt string `json:"created_at"`
}

// RequestPagination - список заявок с пагинацией
type RequestPagination struct {
	Items      []RequestResponse `json:"items"`
	TotalItems int               `json:"total_items"`
	TotalPages int               `json:"total_pages"`
}

// ResponseRequestPagination - ответ со списком заявок
type ResponseRequestPagination struct {
	BaseResponse
	Value *RequestPagination `json:"value"`
}

// CreateGuildRequest - запрос на создание гильдии
type CreateGuildRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
}

// EditMemberRequest - запрос на изменение участника
type EditMemberRequest struct {
	RoleID   int    `json:"role_id"`   // id новой роли
	UserName string `json:"user_name"` // Новое имя (опционально)
}

// EditGuildRequest - запрос на изменение гильдии
type EditGuildRequest struct {
	Title       string `json:"title"`       // Новое название гильдии
	Description string `json:"description"` // Новое описание гильдии
}

// WarStatus - статус войны гильдий
type WarStatus string

const (
	WarStatusPending  WarStatus = "pending"
	WarStatusActive   WarStatus = "active"
	WarStatusFinished WarStatus = "finished"
	WarStatusDeclined WarStatus = "declined"
	WarStatusCanceled WarStatus = "canceled"
	WarStatusExpired  WarStatus = "expired"
)

// DeclareWarRequest - запрос на объявление войны
type DeclareWarRequest struct {
	InitiatorGuildID int `json:"initiator_guild_id"`
	TargetGuildID    int `json:"target_guild_id"`
	InitiatorOwnerID int `json:"initiator_owner_id"`
}

// DeclareWarResponse - ответ на объявление войны
type DeclareWarResponse struct {
	WarID            int       `json:"war_id"`
	InitiatorGuildID int       `json:"initiator_guild_id"`
	TargetGuildID    int       `json:"target_guild_id"`
	Status           WarStatus `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

// ConfirmWarRequest - запрос на подтверждение участия в войне
type ConfirmWarRequest struct {
	TargetOwnerID int `json:"target_owner_id"`
}

// ConfirmWarResponse - ответ на подтверждение участия в войне
type ConfirmWarResponse struct {
	WarID            int       `json:"war_id"`
	InitiatorGuildID int       `json:"initiator_guild_id"`
	TargetGuildID    int       `json:"target_guild_id"`
	Status           WarStatus `json:"status"`
	UpdatedAt        time.Time `json:"updated_at"`
	InitiatorOwnerID int       `json:"initiator_owner_id"`
	TargetOwnerID    int       `json:"target_owner_id"`
}

// CancelWarRequest - запрос на отмену участия в войне
type CancelWarRequest struct {
	OwnerID int `json:"owner_id"`
}

// CancelWarResponse - ответ на отмену участия в войне
type CancelWarResponse struct {
	WarID            int       `json:"war_id"`
	Status           WarStatus `json:"status"`
	CancelledBy      int       `json:"cancelled_by"`
	CancelledAt      time.Time `json:"cancelled_at"`
	InitiatorGuildID int       `json:"initiator_guild_id"`
	TargetGuildID    int       `json:"target_guild_id"`
	InitiatorOwnerID int       `json:"initiator_owner_id"`
	TargetOwnerID    int       `json:"target_owner_id"`
}

// GuildWarItem - элементы списка войн гильдий
type GuildWarItem struct {
	ID               int       `json:"id"`
	InitiatorGuildID int       `json:"initiator_guild_id"`
	TargetGuildID    int       `json:"target_guild_id"`
	Status           WarStatus `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

// GuildWarListResponse - ответ со списком войн гильдий
type GuildWarListResponse struct {
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	Total      int            `json:"total"`
	TotalPages int            `json:"total_pages"`
	Results    []GuildWarItem `json:"results"`
}
