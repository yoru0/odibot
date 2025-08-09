package models

import (
	"math/rand"
	"sort"
)

type Deck []Card

func NewDeck() Deck {
	d := make(Deck, 0, 52)
	for s := Diamonds; s <= Spades; s++ {
		for r := Three; r <= Two; r++ {
			d = append(d, Card{Suit: s, Rank: r})
		}
	}
	return d
}

func (d Deck) Shuffle() {
	rand.Shuffle(len(d), func(i, j int) { d[i], d[j] = d[j], d[i] })
}

func (d Deck) SortRankAsc() {
	sort.Slice(d, func(i, j int) bool {
		if d[i].Rank != d[j].Rank {
			return d[i].Rank < d[j].Rank
		}
		return d[i].Suit < d[j].Suit
	})
}

func (d Deck) SortRankDsc() {
	sort.Slice(d, func(i, j int) bool {
		if d[i].Rank != d[j].Rank {
			return d[i].Rank > d[j].Rank
		}
		return d[i].Suit < d[j].Suit
	})
}

func (d Deck) SortSuitAsc() {
	sort.Slice(d, func(i, j int) bool {
		if d[i].Suit != d[j].Suit {
			return d[i].Suit < d[j].Suit
		}
		return d[i].Rank < d[j].Rank
	})
}

func (d Deck) SortSuitDsc() {
	sort.Slice(d, func(i, j int) bool {
		if d[i].Suit != d[j].Suit {
			return d[i].Suit > d[j].Suit
		}
		return d[i].Rank < d[j].Rank
	})
}

func (d *Deck) Deal(players int) []Deck {
	if players <= 0 {
		return nil
	}
	per := len(*d) / players
	hands := make([]Deck, players)
	for i := range players {
		start := i * per
		end := start + per
		hands[i] = (*d)[start:end]
	}
	*d = (*d)[players*per:]
	return hands
}

func (d *Deck) DealThreePlayer() []Deck {
	const per = 13
	if len(*d) < per*3 {
		return nil
	}
	hands := make([]Deck, 3)
	for i := range 3 {
		start := i * per
		end := start + per
		hands[i] = (*d)[start:end]
	}
	*d = (*d)[per*3:]
	return hands
}

func (d *Deck) RemoveThreeDiamonds() {
	for i, card := range *d {
		if card.Rank == 3 && card.Suit == Diamonds {
			*d = append((*d)[:i], (*d)[i+1:]...)
			return
		}
	}
}
