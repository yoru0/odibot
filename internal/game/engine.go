package game

import (
	"errors"
	"fmt"
)

// Start starts the game.
func (g *Game) Start() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.started {
		return errors.New("already started")
	}
	if len(g.players) < 3 {
		return errors.New("need at least 3 players")
	}

	d := NewDeck()
	d.Shuffle()
	hands := d.Deal(len(g.players))
	for i, player := range g.players {
		player.Hand = hands[i]
		player.Skipped = false
		player.Finished = false
	}
	g.handSize = len(hands[0])

	start := 0
	foundStarter := false
	threeSuits := []Suit{Spades, Hearts, Clubs, Diamonds}
	for _, suit := range threeSuits {
		threeCard := Card{Rank: R3, Suit: suit}
		for i, p := range g.players {
			if p.HasCard(threeCard) >= 0 {
				start = i
				foundStarter = true
				break
			}
		}
		if foundStarter {
			break
		}
	}

	g.turn = start
	g.lead = start
	g.current = Combo{Type: ComboNone}
	g.skipsInRow = 0
	g.started = true
	return nil
}

// Play processes a player's move (playing cards).
func (g *Game) Play(userID string, codes []string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if !g.started {
		return "", errors.New("game not started")
	}

	cp := g.players[g.turn]
	if cp.UserID != userID {
		return "", errors.New("not your turn")
	}
	if len(codes) == 0 {
		return "", errors.New("provide cards to play, e.g.`play 2H 2S`")
	}

	if g.current.Type != ComboNone && cp.Skipped {
		return "", errors.New("you already skipped this round; wait until the table resets")
	}

	selected := make([]Card, 0, len(codes))
	for _, code := range codes {
		c, ok := ParseCard(code)
		if !ok {
			return "", fmt.Errorf("invalid card: %s", code)
		}
		if cp.HasCard(c) == -1 {
			return "", fmt.Errorf("you don't have %s", c.String())
		}
		selected = append(selected, c)
	}

	combo, err := EvaluateCombo(selected)
	if err != nil {
		return "", err
	}
	if g.current.Type != ComboNone && !Beats(combo, g.current) {
		return "", errors.New("your combo does not beat the table")
	}

	if g.current.Type == ComboNone {
		g.resetAllSkips()
	}

	cp.RemoveCards(combo.Cards)
	g.current = combo
	g.lead = g.turn
	g.skipsInRow = 0
	cp.Skipped = false

	if len(cp.Hand) == 0 {
		cp.Finished = true
		g.finishedOrder = append(g.finishedOrder, cp.Idx)
		if g.activePlayers() == 1 {
			g.gameOver = true
			if len(g.finishedOrder) > 0 {
				g.winnerIdx = g.finishedOrder[0]
			}
			return fmt.Sprintf("%s plays %s and finishes. Game over.", cp.Name, comboString(combo)), nil
		}
	}

	g.advanceTurn()
	return fmt.Sprintf("%s plays %s. Next: %s", cp.Name, comboString(combo), g.players[g.turn].Name), nil
}

// Skip handles a player skipping their turn.
func (g *Game) Skip(userID string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if !g.started {
		return "", errors.New("game not started")
	}
	cp := g.players[g.turn]
	if cp.UserID != userID {
		return "", errors.New("not your turn")
	}
	if g.current.Type == ComboNone {
		return "", errors.New("cannot skip on an empty table")
	}

	cp.Skipped = true
	g.skipsInRow++
	g.advanceTurn()

	if g.skipsInRow >= g.activePlayers()-1 {
		g.current = Combo{Type: ComboNone}
		g.skipsInRow = 0
		g.resetAllSkips()
		g.turn = g.lead
		g.advanceTurn()
		leader := g.players[g.turn].Name
		return fmt.Sprintf("%s skips. Table cleared. %s to lead.", cp.Name, leader), nil
	}
	return fmt.Sprintf("%s skips. Next: %s", cp.Name, g.players[g.turn].Name), nil
}
