package combo

import "github.com/yoru0/odibot/capsa/deck"

func CheckStrongerCombo(comboToBeat Combo, currCombo Combo) bool {
	comboToBeatLength := len(comboToBeat.Cards)
	currComboLength := len(currCombo.Cards)

	// If card length not the same.
	if comboToBeatLength != currComboLength {
		return false
	}

	switch comboToBeatLength {
	case 0:
		return true
	case 1, 2, 3:
		return comparePower(currCombo.Power, comboToBeat.Power)
	case 5:
		if currCombo.Type == comboToBeat.Type {
			return comparePower(currCombo.Power, comboToBeat.Power)
		}
		return currCombo.Type > comboToBeat.Type
	}
	return false
}

func comparePower(a, b deck.Card) bool {
	if a.Rank > b.Rank {
		return true
	}
	if a.Rank == b.Rank && a.Suit > b.Suit {
		return true
	}
	return false
}
