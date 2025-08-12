package game

import (
	"fmt"
	"strings"
)

func clone(cards []Card) []Card {
	temp := make([]Card, len(cards))
	copy(temp, cards)
	return temp
}

func discardThrees(in []Card) (out []Card, had3S bool) {
	out = out[:0]
	for _, c := range in {
		if c.Rank == R3 {
			if c.Suit == Spades {
				had3S = true
			}
			continue
		}
		out = append(out, c)
	}
	return out, had3S
}

func comboString(c Combo) string {
	var parts []string
	for _, cc := range c.Cards {
		parts = append(parts, cc.String())
	}
	return fmt.Sprintf("%s [%s]", comboTypeName(c.Type), strings.Join(parts, " "))
}

func comboTypeName(t ComboType) string {
	switch t {
	case ComboSingle:
		return "Single"
	case ComboPair:
		return "Pair"
	case ComboStraight:
		return "Straight"
	case ComboFlush:
		return "Flush"
	case ComboFullHouse:
		return "FullHouse"
	case ComboFourKind:
		return "FourKind"
	case ComboStraightFlush:
		return "StraightFlush"
	default:
		return "None"
	}
}
