package game

import (
	"github.com/yoru0/odibot/capsa/combo"
	"github.com/yoru0/odibot/capsa/deck"
	"github.com/yoru0/odibot/capsa/player"
)

type Game struct {
	Players       player.Players
	CurrIndex     int
	Round         int
	PlayerSkipped int
	ComboPlayed   combo.Combo
	ComboToBeat   combo.Combo
	PlayHistory   deck.CardHistory
}

func NewGame(numPlayers int, userID, username string) *Game {
	p := player.NewPlayers(numPlayers, userID, username)
	first := p.WhoPlaysFirst()
	p.RemoveThree()
	return &Game{
		Players:       p,
		CurrIndex:     first,
		Round:         1,
		PlayerSkipped: 0,
		ComboPlayed:   combo.Combo{Type: combo.None},
		ComboToBeat:   combo.Combo{Type: combo.None},
		PlayHistory:   nil,
	}
}

func (g *Game) ResetPlayerSkips() {
	for i := range g.Players {
		g.Players[i].Skip = false
	}
	g.PlayerSkipped = 0
}

func (g *Game) ResetLastCombo() {
	g.ComboToBeat.Type = combo.None
}

func (g *Game) NextPlayerTurn() {
	g.CurrIndex = (g.CurrIndex + 1) % len(g.Players)
}

func (g *Game) SetCurrentPlayer(name string) {
	for i := range g.Players {
		if g.Players[i].Username == name {
			g.CurrIndex = g.Players[i].ID
			return
		}
	}
}
