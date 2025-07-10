package scoreboard

// UserStat - статистика пользователя
type UserStat struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	Gold                  int    `json:"gold"`
	GoldRatingPos         int    `json:"gold_rating_pos"` // рейтинг по золоту
	Experience            int    `json:"experience"`
	ExpRatingPos          int    `json:"exp_rating_pos"` // рейтинг по опыту
	Rating                int    `json:"rating"`
	RatingRatingPos       int    `json:"rating_rating_pos"` // общий рейтинге
	ChestsOpened          int    `json:"chest_opened"`      // количество открытых сундуков
	ChestsOpenedRatingPos int    `json:"chest_opened_pos"`  // рейтинг по открытым сундукам
}

// GuildStat - статистика гильдии
type GuildStat struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	GuildMembers           int    `json:"guild_members"`             // количество участников в гильдии
	GuildMembersRatingPos  int    `json:"guild_members_rating_pos"`  // рейтинг по количеству участников
	WarsVictories          int    `json:"wars_victories"`            // победы в войнах гильдий
	WarsVictoriesRatingPos int    `json:"wars_victories_rating_pos"` // рейтинг по победам
}

// UserListResponse - ответ со списком пользователей
type UserListResponse struct {
	Page       int        `json:"page"`        // текущая страница
	PageAmount int        `json:"page_amount"` // всего страниц
	Items      []UserStat `json:"items"`       // список пользователей
}

// GuildListResponse - ответ со списком гильдий
type GuildListResponse struct {
	Page       int         `json:"page"`        // текущая страница
	PageAmount int         `json:"page_amount"` // всего страниц
	Items      []GuildStat `json:"items"`       // список гильдий
}
