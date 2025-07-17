package matchmaking

import "net/url"

var matchmakingWsUrl = url.URL{Scheme: "ws", Host: "37.9.53.32:80"}

type MatchmakingPath *url.URL

var (
	RandomMatchmaking = MatchmakingPath(&url.URL{Path: "/matchmaking/random"})
	RankedMatchmaking = MatchmakingPath(&url.URL{Path: "/matchmaking/ranked"})
	CustomMatchmaking = MatchmakingPath(&url.URL{Path: "/matchmaking/custom"})
)

func formatMatchmakingPath(path MatchmakingPath) string {
	url := matchmakingWsUrl.ResolveReference(path)

	return url.String()
}
