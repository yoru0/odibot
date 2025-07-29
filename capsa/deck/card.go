package deck

import "fmt"

type Rank int

const (
	Three Rank = iota + 3
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
	Two
)

func (r Rank) String() string {
	switch r {
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	case Two:
		return "2"
	default:
		return fmt.Sprint(int(r))
	}
}

type Suit int

const (
	Diamonds Suit = iota
	Clubs
	Hearts
	Spades
)

func (s Suit) String() string {
	switch s {
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	case Hearts:
		return "♥"
	case Spades:
		return "♠"
	}
	return "?"
}

type Card struct {
	Rank Rank
	Suit Suit
}
