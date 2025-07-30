package player

import "github.com/yoru0/odibot/capsa/deck"

type Player struct {
	ID         int
	UserID     string
	Username   string
	GlobalName string
	Joined     bool
	Hand       deck.Deck
	Skip       bool
}

type Players []Player
