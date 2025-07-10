package handlers

type PlayerStats struct {
	Username   string `json:"username"`
	Wins       int    `json:"wins"`
	Losses     int    `json:"losses"`
	TotalGames int    `json:"total_games"`
	Rank       int    `json:"rank"`
}

type GuildsMemberStats struct {
	Username        string `json:"username"`
	Wins            int    `json:"wins"`
	GuildChests     int    `json:"guildChests"`
	GuildRagePoints int    `json:"guild_rage_points"`
	WarWins         int    `json:"war_wins"`
	TotalDonations  int    `json:"total_donations"`
}

type GuildInternalStats struct {
	GuildName string              `json:"guild_name"`
	GuildTag  string              `json:"guild_tag"`
	Members   []GuildsMemberStats `json:"members"`
	TotalRage int                 `json:"total_rage"`
}

func MyStatsHandler(token string) (PlayerStats, error) {
	return PlayerStats{
		Username: "current_user", Wins: 15, Losses: 5, TotalGames: 20, Rank: 42}, nil
}

func PlayersStatsHandler() ([]PlayerStats, error) {
	return []PlayerStats{
		{Rank: 1, Username: "naval_legend", Wins: 200, Losses: 20, TotalGames: 220},
		{Rank: 2, Username: "sea_wolf", Wins: 180, Losses: 25, TotalGames: 205},
		{Rank: 3, Username: "admiral", Wins: 150, Losses: 30, TotalGames: 180},
		{Rank: 4, Username: "captain", Wins: 120, Losses: 35, TotalGames: 155},
		{Rank: 5, Username: "first_mate", Wins: 100, Losses: 40, TotalGames: 140},
		{Rank: 6, Username: "boatswain", Wins: 90, Losses: 45, TotalGames: 135},
		{Rank: 7, Username: "quartermaster", Wins: 80, Losses: 50, TotalGames: 130},
		{Rank: 8, Username: "carpenter", Wins: 70, Losses: 55, TotalGames: 125},
		{Rank: 9, Username: "gunner", Wins: 60, Losses: 60, TotalGames: 120},
		{Rank: 10, Username: "cook", Wins: 50, Losses: 70, TotalGames: 120},
	}, nil
}

func GuildInternalStatsHandler(token string) (GuildInternalStats, error) {
	return GuildInternalStats{
		GuildName: "Морские волки",
		GuildTag:  "WOLF",
		TotalRage: 3250,
		Members: []GuildsMemberStats{
			{Username: "admiral", Wins: 150, GuildChests: 25, GuildRagePoints: 750, WarWins: 15, TotalDonations: 1200},
			{Username: "captain", Wins: 120, GuildChests: 20, GuildRagePoints: 600, WarWins: 12, TotalDonations: 900},
			{Username: "first_mate", Wins: 100, GuildChests: 18, GuildRagePoints: 500, WarWins: 10, TotalDonations: 750},
			{Username: "navigator", Wins: 90, GuildChests: 15, GuildRagePoints: 450, WarWins: 8, TotalDonations: 600},
			{Username: "boatswain", Wins: 80, GuildChests: 12, GuildRagePoints: 400, WarWins: 7, TotalDonations: 500},
			{Username: "quartermaster", Wins: 70, GuildChests: 10, GuildRagePoints: 350, WarWins: 6, TotalDonations: 400},
			{Username: "gunner", Wins: 60, GuildChests: 8, GuildRagePoints: 300, WarWins: 5, TotalDonations: 300},
		},
	}, nil
}
