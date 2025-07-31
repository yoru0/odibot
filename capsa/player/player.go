package player

import (
	"github.com/yoru0/odibot/capsa/deck"
	"github.com/yoru0/odibot/capsa/design"
)

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

func NewPlayers(numPlayers int, userID, username string) Players {
	d := deck.NewDeck()
	d.ShuffleDeck()
	hands := d.DealDeck(numPlayers)
	players := make(Players, numPlayers)

	for i := range hands {
		hands[i].SortDeckByRankAscending()
	}

	for i := range numPlayers {
		players[i] = Player{
			ID:       i,
			UserID:   userID,
			Username: username,
			Hand:     hands[i],
			Skip:     false,
		}
	}
	return players
}

func (p Player) ShowPlayerHand() {
	for i, card := range p.Hand {
		design.PrintPlayersHandWithColor(card, i)
	}
}

func (p *Players) RemovePlayerAfterWin(id int) {
	if id < 0 || id >= len(*p) {
		return
	}
	*p = append((*p)[:id], (*p)[id+1:]...)
}

func (p *Player) SkipTurn() {
	p.Skip = true
}
