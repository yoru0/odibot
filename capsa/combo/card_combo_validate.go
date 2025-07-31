package combo

import (
	"fmt"

	"github.com/yoru0/odibot/capsa/deck"
)

func CheckCombo(cards deck.Deck) (bool, Combo) {
	valid := false
	var combo Combo
	cards.SortDeckByRankAscending()

	length := len(cards)

	switch length {
	case 0:
		valid = true
		combo = Combo{
			Type:  Skip,
			Cards: cards,
		}
	case 1:
		valid = true
		combo = Combo{
			Type:  Single,
			Cards: cards,
			Power: cards[0],
		}
	case 2:
		if IsPair(cards) {
			valid = true
			combo = Combo{
				Type:  Pair,
				Cards: cards,
				Power: cards[1],
			}
		}
	case 3:
		if IsTriple(cards) {
			valid = true
			combo = Combo{
				Type:  Triple,
				Cards: cards,
				Power: cards[2],
			}
		}
	case 5:
		switch {
		case IsStraightFlush(cards):
			valid = true
			combo = Combo{
				Type:  StraightFlush,
				Cards: cards,
				Power: cards[4],
			}
		case IsFourOfAKind(cards):
			valid = true
			combo = Combo{
				Type:  FourOfAKind,
				Cards: cards,
				Power: cards[3],
			}
		case IsFullHouse(cards):
			valid = true
			combo = Combo{
				Type:  FullHouse,
				Cards: cards,
				Power: cards[2],
			}
		case IsFlush(cards):
			valid = true
			combo = Combo{
				Type:  Flush,
				Cards: cards,
				Power: cards[4],
			}
		case IsStraight(cards):
			valid = true
			combo = Combo{
				Type:  Straight,
				Cards: cards,
				Power: cards[4],
			}
		default:
			valid = false
			combo = Combo{
				Type: None,
			}
			fmt.Println("Invalid 5-card combo.")
		}
	default:
		valid = false
		combo = Combo{
			Type: None,
		}
		fmt.Println("Invalid combo length.")
	}
	return valid, combo
}
