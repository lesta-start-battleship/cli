package models

import (
	"lesta-start-battleship/cli/internal/api/guilds"
)

type OAuthPollingResultMsg struct {
	ID       int
	Username string
	Gold     int
	Error    string
}

type AuthSuccessMsg struct {
	ID       int
	Username string
	Gold     int
}

type LogoutMsg struct{}

type OpenChatMsg struct {
	GuildID int
}

type UsernameChangeMsg struct {
	NewUsername string
	Gold        int
}

type ChatKeyHandledMsg struct{}

type ChatClosedMsg struct{}

type GuildDataMsg struct {
	Member *guilds.MemberResponse
	Guild  *guilds.GuildResponse
}

type GuildNoMemberMsg struct{}

type MemberRoleChangeMsg struct {
	Username string
}

type MemberDeleteMsg struct {
	Username string
}

type GuildExitedMsg struct{}

type RequestProcessedMsg struct {
	Message string
}

type GuildCreatedMsg struct {
	Guild *guilds.GuildResponse
}

type JoinRequestSentMsg struct{}

type GuildEditMsg struct {
	Guild *guilds.GuildResponse
}

type GuildDeletedMsg struct{}

type DeclareWarMsg struct{}

type WarRequestProcessedMsg struct {
	Message string
}
