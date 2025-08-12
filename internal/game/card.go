package game

import (
	"fmt"
	"sort"
	"strings"
)

type Suit int
type Rank int

const (
	Diamonds Suit = iota
	Clubs
	Hearts
	Spades
)

const (
	R3  Rank = 3
	R4  Rank = 4
	R5  Rank = 5
	R6  Rank = 6
	R7  Rank = 7
	R8  Rank = 8
	R9  Rank = 9
	R10 Rank = 10
	J   Rank = 11
	Q   Rank = 12
	K   Rank = 13
	A   Rank = 14
	R2  Rank = 15
)

type Card struct {
	Suit Suit
	Rank Rank
}

func (c Card) String() string {
	return rankString(c.Rank) + suitString(c.Suit)
}

func rankString(r Rank) string {
	switch r {
	case R10:
		return "10"
	case J:
		return "J"
	case Q:
		return "Q"
	case K:
		return "K"
	case A:
		return "A"
	case R2:
		return "2"
	default:
		return fmt.Sprintf("%d", int(r))
	}
}

func suitString(s Suit) string {
	switch s {
	case Diamonds:
		return "D"
	case Clubs:
		return "C"
	case Hearts:
		return "H"
	default:
		return "S"
	}
}

func ParseCard(code string) (Card, bool) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return Card{}, false
	}
	var suit Suit
	switch code[len(code)-1] {
	case 'D':
		suit = Diamonds
	case 'C':
		suit = Clubs
	case 'H':
		suit = Hearts
	case 'S':
		suit = Spades
	default:
		return Card{}, false
	}
	rankPart := code[:len(code)-1]
	var r Rank
	switch rankPart {
	case "J":
		r = J
	case "Q":
		r = Q
	case "K":
		r = K
	case "A":
		r = A
	case "2":
		r = R2
	case "10":
		r = R10
	default:
		if len(rankPart) == 1 && rankPart[0] >= '3' && rankPart[0] <= '9' {
			r = Rank(rankPart[0] - '0')
		} else {
			return Card{}, false
		}
	}
	return Card{Rank: r, Suit: suit}, true
}

// SortCardsAsc sorts Rank, then Suit ascending.
func SortCardsAsc(cards []Card) {
	sort.Slice(cards, func(i, j int) bool {
		if cards[i].Rank == cards[j].Rank {
			return cards[i].Suit < cards[j].Suit
		}
		return cards[i].Rank < cards[j].Rank
	})
}

// SortCardsDesc sorts Rank, then Suit descending.
func SortCardsDesc(cards []Card) {
	sort.Slice(cards, func(i, j int) bool {
		if cards[i].Rank == cards[j].Rank {
			return cards[i].Suit > cards[j].Suit
		}
		return cards[i].Rank > cards[j].Rank
	})
}
