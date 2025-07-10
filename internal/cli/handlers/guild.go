package handlers

type GuildResponse struct {
	Member bool      `json:"member"`
	Owner  bool      `json:"owner"`
	Info   GuildInfo `json:"info"`
}

type GuildInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

func GetGuildInfo(token string) (GuildResponse, error) {
	return GuildResponse{
		Member: true,
		Owner:  false,
		Info: GuildInfo{
			Id:   1,
			Name: "Sea Wolf",
			Tag:  "WOLF",
		},
	}, nil
}
