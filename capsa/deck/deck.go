package deck

import (
	"math/rand"
	"sort"
)

type Deck []Card

func NewDeck() Deck {
	var d Deck
	for i := Diamonds; i <= Spades; i++ {
		for j := Three; j <= Two; j++ {
			d = append(d, Card{Suit: i, Rank: j})
		}
	}
	return d
}

func (d Deck) ShuffleDeck() {
	rand.Shuffle(len(d), func(i, j int) { d[i], d[j] = d[j], d[i] })
}

func (d Deck) SortDeckByRankAscending() {
	sort.Slice(d, func(i, j int) bool {
		if d[i].Rank != d[j].Rank {
			return d[i].Rank < d[j].Rank
		}
		return d[i].Suit < d[j].Suit
	})
}

func (d Deck) SortDeckByRankDescending() {
	sort.Slice(d, func(i, j int) bool {
		if d[i].Rank != d[j].Rank {
			return d[i].Rank > d[j].Rank
		}
		return d[i].Suit < d[j].Suit
	})
}

func (d Deck) SortDeckBySuitAscending() {
	sort.Slice(d, func(i, j int) bool {
		if d[i].Suit != d[j].Suit {
			return d[i].Suit < d[j].Suit
		}
		return d[i].Rank < d[j].Rank
	})
}

func (d Deck) SortDeckBySuitDescending() {
	sort.Slice(d, func(i, j int) bool {
		if d[i].Suit != d[j].Suit {
			return d[i].Suit > d[j].Suit
		}
		return d[i].Rank < d[j].Rank
	})
}

func (d *Deck) DealDeck(players int) []Deck {
	per := len(*d) / players
	hands := make([]Deck, players)
	for i := range players {
		start := i * per
		end := start + per
		hands[i] = (*d)[start:end]
	}
	total := players * per
	*d = (*d)[total:]
	return hands
}



func (d *Deck) RemoveThreeDiamonds() {
	for i, card := range *d {
		if  card.Rank == 3 && card.Suit == Diamonds {
			*d = append((*d)[:i], (*d)[i+1:]...)
			return
		}
	}
}
