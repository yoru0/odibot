package combo

import (
	"sort"

	"github.com/yoru0/odibot/internal/models"
)

func IsPair(cards models.Deck) bool {
	if len(cards) != 2 {
		return false
	}
	return cards[0].Rank == cards[1].Rank
}

func IsTriple(cards models.Deck) bool {
	if len(cards) != 3 {
		return false
	}
	return cards[0].Rank == cards[1].Rank && cards[1].Rank == cards[2].Rank
}

func IsStraight(cards models.Deck) bool {
	if len(cards) != 5 {
		return false
	}

	var rank []int
	for _, c := range cards {
		rank = append(rank, int(c.Rank))
	}

	sort.Ints(rank)
	for i := 0; i < len(rank)-1; i++ {
		if rank[i+1] != rank[i]+1 {
			return false
		}
	}
	return true
}

func IsFlush(cards models.Deck) bool {
	if len(cards) != 5 {
		return false
	}

	for i := 0; i < len(cards)-1; i++ {
		if cards[i].Suit != cards[i+1].Suit {
			return false
		}
	}
	return true
}

func IsFullHouse(cards models.Deck) bool {
	if len(cards) != 5 {
		return false
	}

	var rank []int
	for _, c := range cards {
		rank = append(rank, int(c.Rank))
	}

	sort.Ints(rank)
	if rank[0] == rank[2] && rank[3] == rank[4] && rank[2] != rank[3] {
		return true
	}
	if rank[0] == rank[1] && rank[2] == rank[4] && rank[1] != rank[2] {
		return true
	}
	return false
}

func IsFourOfAKind(cards models.Deck) bool {
	if len(cards) != 5 {
		return false
	}

	var rank []int
	for _, c := range cards {
		rank = append(rank, int(c.Rank))
	}

	sort.Ints(rank)
	if rank[0] == rank[3] || rank[1] == rank[4] {
		return true
	}
	return false
}

func IsStraightFlush(cards models.Deck) bool {
	return IsStraight(cards) && IsFlush(cards)
}
