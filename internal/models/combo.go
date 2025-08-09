package models

type ComboType int

const (
	None ComboType = iota
	Skip
	Single
	Pair
	Triple
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	InvalidCombo
)

type Combo struct {
	Type  ComboType
	Cards Deck
	Power Card
}

func (c ComboType) String() string {
	switch c {
	case None:
		return "None"
	case Skip:
		return "Skip"
	case Single:
		return "Single"
	case Pair:
		return "Pair"
	case Triple:
		return "Triple"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "Full House"
	case FourOfAKind:
		return "Four of a Kind"
	case StraightFlush:
		return "Straight Flush"
	case InvalidCombo:
		return "Invalid Combo"
	default:
		return "Unknown"
	}
}
