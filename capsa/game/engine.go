package game

import (
	"fmt"

	"github.com/yoru0/odibot/capsa/combo"
	"github.com/yoru0/odibot/capsa/deck"
	"github.com/yoru0/odibot/capsa/design"
	"github.com/yoru0/odibot/capsa/player"
	"github.com/yoru0/odibot/capsa/winner"
)

func StartGame() {
	g := NewGame(4, "p", "p")
	var w winner.Winner
	var lastPlayerThatPlay string

	for len(g.Players)-len(w) > 1 {
		if g.Round != 1 {
			g.SetCurrentPlayer(lastPlayerThatPlay)
		}
		var sp player.SkipedPlayerName

		for g.PlayerSkipped < (len(g.Players) - len(w) - 1) {
			for g.Players[g.CurrIndex].Skip || len(g.Players[g.CurrIndex].Hand) == 0 {
				g.NextPlayerTurn()
			}

			detailBar(*g, sp)

			fmt.Printf("%s [%d] ~\n", g.Players[g.CurrIndex].Username, len(g.Players[g.CurrIndex].Hand))
			g.Players[g.CurrIndex].ShowPlayerHand()

			var cards deck.Deck
			g.ComboPlayed, cards = g.Players[g.CurrIndex].PickCard(g.ComboToBeat)

			g.PlayHistory = append(g.PlayHistory, deck.History{
				CardsPlayed: cards,
				PlayerName:  g.Players[g.CurrIndex].Username,
			})

			if g.ComboPlayed.Type == combo.Skip {
				g.Players[g.CurrIndex].Skip = true
				g.PlayerSkipped++
				g.PlayHistory.RemoveNilSkipHistory()
				sp = append(sp, g.Players[g.CurrIndex].Username)
				fmt.Printf("You skipped.\n\n")
			} else {
				g.ComboToBeat = g.ComboPlayed
				fmt.Print("You played: ")
				for i := range g.ComboPlayed.Cards {
					design.PrintIndividualCardWithColor(cards[i])
				}
				fmt.Printf("[%s]\n\n", g.ComboPlayed.Type)
			}

			if len(g.Players[g.CurrIndex].Hand) == 0 {
				w.AppendWinner(g.Players[g.CurrIndex].Username)
				fmt.Printf("You win! Position: %d\n\n", len(w))
				resetToNewRound(g)
			}

			// g.CheckGameDetails(g.CurrIndex)

			design.PressEnterToContinue()

			g.NextPlayerTurn()
		}

		if len(g.PlayHistory) > 0 {
			lastPlayerThatPlay = g.PlayHistory[len(g.PlayHistory)-1].PlayerName
		}
		resetToNewRound(g)
	}
	w.ShowWinners()
}

func resetToNewRound(g *Game) {
	g.Round++
	g.ResetPlayerSkips()
	g.ResetLastCombo()
	g.PlayHistory.ResetHistory()
}

func detailBar(g Game, sp player.SkipedPlayerName) {
	fmt.Println("──────────────────────────────────────────────────────────────────────────────────────────────────")
	fmt.Printf("Round: %-15d Skips: %d ", g.Round, g.PlayerSkipped)

	if len(sp) > 0 {
		fmt.Print("[")
		for i := range sp {
			if i == len(sp)-1 {
				fmt.Printf("%s]", sp[i])
			} else {
				fmt.Printf("%s, ", sp[i])
			}
		}
	}

	fmt.Println()
	if len(g.PlayHistory) > 0 {
		g.PlayHistory.ShowHistory()
	}

	fmt.Printf("──────────────────────────────────────────────────────────────────────────────────────────────────\n\n")
}
