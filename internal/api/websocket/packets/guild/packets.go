package guild

type Packet interface {
	isGuildPacket()
}

type ChatMessage struct {
	Msg string `json:"content"`
}

func (ChatMessage) isGuildPacket() {}

type ChatHistory struct {
	Type string               `json:"type"`
	Data []ChatHistoryMessage `json:"data"`
}

func (ChatHistory) isGuildPacket() {}

type ChatHistoryMessage struct {
	Id        string `json:"_id"`
	GuildId   int    `json:"guild_id"`
	UserId    int    `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

func (ChatHistoryMessage) isGuildPacket() {}

type Disconnect struct{}

func (Disconnect) isGuildPacket() {}
