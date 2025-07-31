package player

import (
	"fmt"

	"github.com/yoru0/odibot/capsa/deck"
	"github.com/yoru0/odibot/capsa/design"
)

func (p *Players) WhoPlaysFirst() int {
	for i := range *p {
		for _, card := range (*p)[i].Hand {
			if card.Rank == deck.Three && card.Suit == deck.Spades {
				return i
			}
		}
	}
	return -1
}

func (p *Players) RemoveThree() {
	fmt.Println("Threes in hand ~")
	for i := range *p {
		var picks []int
		hasThreeSpades := false
		for _, card := range (*p)[i].Hand {
			if card.Rank == deck.Three && card.Suit == deck.Spades {
				hasThreeSpades = true
				break
			}
		}
		for j, card := range (*p)[i].Hand {
			if card.Rank == deck.Three {
				if len(picks) == 0 {
					if hasThreeSpades {
						fmt.Printf("->  %s -- ", (*p)[i].Username)
					} else {
						fmt.Printf("    %s -- ", (*p)[i].Username)
					}
				}
				design.PrintIndividualCardWithColor(card)
				picks = append(picks, j)
			}
		}
		if len(picks) > 0 {
			(*p)[i].RemovePlayedCards(picks)
			fmt.Println()
		}
	}
	fmt.Println()
}
