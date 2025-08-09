package player

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/yoru0/odibot/capsa/combo"
	"github.com/yoru0/odibot/capsa/deck"
	"github.com/yoru0/odibot/capsa/design"
)

func (p *Player) PickCard(lastCombo combo.Combo) (combo.Combo, deck.Deck) {
	reader := bufio.NewReader(os.Stdin)

	for {
		if lastCombo.Type == combo.None {
			fmt.Println("Choose your cards to play (can't skip this turn).")
			fmt.Print("> ")
		} else {
			fmt.Println("Choose your cards to play or press `Enter` to skip.")
			fmt.Print("> ")
		}
		line, _ := reader.ReadString('\n')
		parts := strings.Fields(line)
		var combos combo.Combo

		picks, cards, valid := p.checkPickedCardIsValid(parts)

		// Check whether the selected cards have a valid combo or not.
		if valid {
			valid, combos = combo.CheckCombo(cards)
			if !valid {
				fmt.Printf("Not a playable combo.\n\n")
				continue
			}
		}

		// First turn on each round can't skip.
		if valid && combos.Type == combo.Skip && lastCombo.Type.String() == "None" {
			valid = false
			fmt.Printf("You can't skip this turn.\n\n")
			continue
		} else if combos.Type == combo.Skip {
			return combos, nil
		}

		// Check whether the combo played is stronger than last combo.
		if valid && lastCombo.Type == combo.None {
			valid = true
		} else if valid {
			valid = combo.CheckStrongerCombo(lastCombo, combos)
			if !valid {
				errorNeedStrongerCards(lastCombo)
				fmt.Println()
				continue
			}
		}

		if valid {
			if combos.Type == combo.Skip {
				return lastCombo, nil
			}
			p.RemovePlayedCards(picks)
			return combos, cards
		}

		fmt.Printf("Please enter valid numbers between 1 and %d.\n\n", len(p.Hand))
	}
}

func (p *Player) RemovePlayedCards(picks []int) {
	sort.Sort(sort.Reverse(sort.IntSlice(picks)))
	for _, i := range picks {
		if i >= 0 && i < len(p.Hand) {
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
		}
	}
}

func (p *Player) checkPickedCardIsValid(parts []string) ([]int, deck.Deck, bool) {
	var picks []int
	var cards deck.Deck
	valid := true
	seen := make(map[int]bool)

	for _, part := range parts {
		num, err := strconv.Atoi(part)

		// Check out of bound case, e.g. [0, 14].
		if err != nil || num < 1 || num > len(p.Hand) {
			fmt.Println("Invalid input:", part)
			valid = false
			break
		}

		// Check duplicates on selected cards, e.g. [1, 1].
		index := num - 1
		if seen[index] {
			fmt.Printf("Duplicate pick: %d\n", num)
			valid = false
			break
		}
		seen[index] = true

		picks = append(picks, index)
		cards = append(cards, p.Hand[index])
	}
	return picks, cards, valid
}

func errorNeedStrongerCards(lastCombo combo.Combo) {
	switch len(lastCombo.Cards) {
	case 1, 2, 3:
		fmt.Printf("Need a `%s` with power higher than ", lastCombo.Type)
		for i := range lastCombo.Cards {
			design.PrintIndividualCardWithColor(lastCombo.Cards[i])
		}
		fmt.Println()
	case 5:
		fmt.Printf("Need a combo stronger than `%s` or cards higher than ", lastCombo.Type)
		for i := range lastCombo.Cards {
			design.PrintIndividualCardWithColor(lastCombo.Cards[i])
		}
		fmt.Println()
	}
}
