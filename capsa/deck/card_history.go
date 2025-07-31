package deck

import "fmt"

type History struct {
	CardsPlayed Deck
	PlayerName  string
}

type CardHistory []History

func (h *CardHistory) ResetHistory() {
	*h = nil
}

func (h CardHistory) ShowHistory() {
	fmt.Println("\nHistory:")
	for i := range h {
		if len(h[i].CardsPlayed) > 0 {
			if i == len(h)-1 {
				fmt.Printf("->  %s -- ", h[i].PlayerName)
			} else {
				fmt.Printf("    %s -- ", h[i].PlayerName)
			}
			for j := range h[i].CardsPlayed {
				printCardWithColor(h[i].CardsPlayed[j])
			}
			fmt.Println()
		}
	}
}

func (h *CardHistory) RemoveNilSkipHistory() {
	*h = (*h)[:len(*h)-1]
}

func printCardWithColor(card Card) {
	colorCode := ""

	switch card.Suit {
	case 0, 2:
		colorCode = "\033[31m" // red
	case 1, 3:
		colorCode = "\033[36m" // cyan
	}
	reset := "\033[0m" // reset

	fmt.Printf("%s%s%s ", colorCode, card.Rank.String()+card.Suit.String(), reset)
}
