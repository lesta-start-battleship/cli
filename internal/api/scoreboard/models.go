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
	ChestsOpened          int    `json:"chests_opened"`     // количество открытых сундуков
	ChestsOpenedRatingPos int    `json:"chests_opened_pos"` // рейтинг по открытым сундукам
}

// GuildStat - статистика гильдии
type GuildStat struct {
	ID                     int    `json:"id"`
	GuildTag               string `json:"tag"`               // тег гильдии
	Players                int    `json:"players"`           // количество участников
	PlayersRatingPos       int    `json:"playes_rating_pos"` // рейтинг участникам (у ребят в Swagger указано поле как будто с ошибкой, но мало ли)
	WarsVictories          int    `json:"wins"`              // победы в войнах гильдий
	WarsVictoriesRatingPos int    `json:"wins_rating_pos"`   // рейтинг по победам
}

// UserListResponse - ответ со списком пользователей
type UserListResponse struct {
	TotalPages int        `json:"total_pages"` // количество страниц
	TotalItems int        `json:"total_items"` // всего значений
	Items      []UserStat `json:"items"`       // список пользователей
}

// GuildListResponse - ответ со списком гильдий
type GuildListResponse struct {
	TotalPages int         `json:"total_pages"` // количество страниц
	TotalItems int         `json:"total_items"` // всего значений
	Items      []GuildStat `json:"items"`       // список гильдий
}
