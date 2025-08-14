package bot

import (
	"fmt"

	"github.com/yoru0/odibot/internal/game"
)

func prettyLabel(c game.Card) string {
	return rankOut(c.Rank) + suitEmoji(c.Suit)
}

func rankOut(r game.Rank) string {
	switch r {
	case game.R10:
		return "10"
	case game.J:
		return "J"
	case game.Q:
		return "Q"
	case game.K:
		return "K"
	case game.A:
		return "A"
	case game.R2:
		return "2"
	default:
		return fmt.Sprintf("%d", int(r))
	}
}

func suitEmoji(s game.Suit) string {
	switch s {
	case game.Diamonds:
		return "♦️"
	case game.Clubs:
		return "♣️"
	case game.Hearts:
		return "♥️"
	default:
		return "♠️"
	}
}
