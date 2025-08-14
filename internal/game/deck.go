package game

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"time"
)

type Deck []Card

func NewDeck() Deck {
	d := make(Deck, 0, 52)
	for s := Diamonds; s <= Spades; s++ {
		for r := R3; r <= R2; r++ {
			d = append(d, Card{Rank: r, Suit: s})
		}
	}
	return d
}

// Shuffle shuffles the deck.
func (d Deck) Shuffle() {
	var seed int64
	binary.Read(cryptoRand.Reader, binary.LittleEndian, &seed)
	r := rand.New(rand.NewSource(time.Now().UnixNano() ^ seed))
	r.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}

// Deal returns a slice of card for each players.
func (d Deck) Deal(players int) [][]Card {
	hands := make([][]Card, players)
	cardsPerPlayer := 13
	for i := range hands {
		hands[i] = make([]Card, 0, cardsPerPlayer)
	}
	for i := 0; i < cardsPerPlayer*players; i++ {
		hands[i%players] = append(hands[i%players], d[i])
	}
	for i := range hands {
		SortCardsAsc(hands[i])
	}
	return hands
}
