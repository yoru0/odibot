package store

import "github.com/yoru0/odibot/internal/game"

type Session struct {
	LobbyChannelID string
	Desired        int
	Game           *game.Game
	DMChannel      map[string]string // userID -> dmChannelID
	Started        bool
	HasDummy       bool

	Selected map[string][]string
}
