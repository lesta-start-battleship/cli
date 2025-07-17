package clientdeps

type PlayerInfo struct {
	id   int
	name string
}

func NewEmptyPlayerInfo() *PlayerInfo {
	return &PlayerInfo{
		id:   0,
		name: "",
	}
}

func NewPlayerInfo(id int, username string) *PlayerInfo {
	return &PlayerInfo{
		id:   id,
		name: username,
	}
}

func (p *PlayerInfo) Id() int {
	return p.id
}

func (p *PlayerInfo) Name() string {
	return p.name
}
