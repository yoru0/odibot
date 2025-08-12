package game

import (
	"errors"
	"fmt"
)

type ComboType int

const (
	ComboNone ComboType = iota
	ComboSingle
	ComboPair
	ComboTriple
	ComboStraight
	ComboFlush
	ComboFullHouse
	ComboFourKind
	ComboStraightFlush
)

type Combo struct {
	Type    ComboType
	Cards   []Card
	KeyRank Rank
	KeySuit Suit
}

// EvaluateCombo returns an available Combo by played cards.
func EvaluateCombo(cards []Card) (Combo, error) {
	n := len(cards)
	switch n {
	case 1:
		return single(cards)
	case 2:
		return pair(cards)
	case 3:
		return triple(cards)
	case 5:
		if c, ok := straightFlush(cards); ok {
			return c, nil
		}
		if c, ok := fourKind(cards); ok {
			return c, nil
		}
		if c, ok := fullHouse(cards); ok {
			return c, nil
		}
		if c, ok := flush(cards); ok {
			return c, nil
		}
		if c, ok := straight(cards); ok {
			return c, nil
		}
	}
	return Combo{}, errors.New("invalid combo")
}

// Beats determines if combo a beats combo b.
func Beats(a, b Combo) bool {
	if b.Type == ComboNone {
		return true
	}
	if len(a.Cards) != len(b.Cards) {
		return false
	}
	if a.Type != b.Type {
		return a.Type > b.Type
	}

	switch a.Type {
	case ComboSingle, ComboPair, ComboTriple, ComboStraight, ComboStraightFlush, ComboFullHouse, ComboFourKind:
		if a.KeyRank != b.KeyRank {
			return a.KeyRank > b.KeyRank
		}
		return a.KeySuit > b.KeySuit
	case ComboFlush:
		if a.KeySuit != b.KeySuit {
			return a.KeySuit > b.KeySuit
		}
		fmt.Println(a.KeyRank, b.KeyRank)
		return a.KeyRank > b.KeyRank
	}
	return false
}

func single(cards []Card) (Combo, error) {
	c := cards[0]
	return Combo{
		Type:    ComboSingle,
		Cards:   []Card{c},
		KeyRank: c.Rank,
		KeySuit: c.Suit,
	}, nil
}

func pair(cards []Card) (Combo, error) {
	if cards[0].Rank != cards[1].Rank {
		return Combo{}, errors.New("pair must have the same rank")
	}
	SortCardsDesc(cards)
	return Combo{
		Type:    ComboPair,
		Cards:   cards,
		KeyRank: cards[0].Rank,
		KeySuit: cards[0].Suit,
	}, nil
}

func triple(cards []Card) (Combo, error) {
	if !(cards[0].Rank == cards[1].Rank && cards[1].Rank == cards[2].Rank) {
		return Combo{}, errors.New("triple must have the same rank")
	}
	SortCardsDesc(cards)
	return Combo{
		Type:    ComboTriple,
		Cards:   cards,
		KeyRank: cards[0].Rank,
		KeySuit: cards[0].Suit,
	}, nil
}

func straight(cards []Card) (Combo, bool) {
	c := clone(cards)
	rank, ok := isStraight(c)
	if !ok {
		return Combo{}, false
	}
	SortCardsDesc(c)
	return Combo{
		Type:    ComboStraight,
		Cards:   c,
		KeyRank: rank,
		KeySuit: c[0].Suit,
	}, true
}

func flush(cards []Card) (Combo, bool) {
	c := clone(cards)
	suit, ok := isFlush(c)
	if !ok {
		return Combo{}, false
	}
	SortCardsDesc(c)
	fmt.Println()
	return Combo{
		Type:    ComboFlush,
		Cards:   c,
		KeyRank: c[0].Rank,
		KeySuit: suit,
	}, true
}

func fullHouse(cards []Card) (Combo, bool) {
	m := countByRank(cards)
	var three Rank
	has3, has2 := false, false
	for rank, count := range m {
		if count == 3 {
			three = rank
			has3 = true
		}
		if count == 2 {
			has2 = true
		}
	}
	if !(has3 && has2) {
		return Combo{}, false
	}
	SortCardsDesc(cards)
	suit := Diamonds
	for _, card := range cards {
		if card.Rank == three && card.Suit > suit {
			suit = card.Suit
		}
	}
	return Combo{
		Type:    ComboFullHouse,
		Cards:   cards,
		KeyRank: three,
		KeySuit: suit,
	}, true
}

func fourKind(cards []Card) (Combo, bool) {
	c := clone(cards)
	m := countByRank(c)
	var four Rank
	found := false
	for rank, count := range m {
		if count == 4 {
			four = rank
			found = true
			break
		}
	}
	if !found {
		return Combo{}, false
	}
	SortCardsDesc(cards)
	suit := Diamonds
	for _, card := range cards {
		if card.Rank == four && card.Suit > suit {
			suit = card.Suit
		}
	}
	return Combo{
		Type:    ComboFourKind,
		Cards:   cards,
		KeyRank: four,
		KeySuit: suit,
	}, true
}

func straightFlush(cards []Card) (Combo, bool) {
	c := clone(cards)
	suit, ok := isFlush(c)
	if !ok {
		return Combo{}, false
	}
	rank, ok := isStraight(c)
	if !ok {
		return Combo{}, false
	}
	return Combo{
		Type:    ComboStraightFlush,
		Cards:   c,
		KeyRank: rank,
		KeySuit: suit,
	}, true
}

func isStraight(cards []Card) (Rank, bool) {
	SortCardsAsc(cards)
	for i := 1; i < len(cards); i++ {
		if cards[i].Rank != cards[i-1].Rank+1 {
			return 0, false
		}
	}
	return cards[len(cards)-1].Rank, true
}

func isFlush(cards []Card) (Suit, bool) {
	s := cards[0].Suit
	for _, c := range cards {
		if c.Suit != s {
			return 0, false
		}
	}
	return s, true
}

func countByRank(cards []Card) map[Rank]int {
	m := make(map[Rank]int)
	for _, c := range cards {
		m[c.Rank]++
	}
	return m
}
