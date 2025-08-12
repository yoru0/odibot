package game

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Game struct {
	mu          sync.Mutex
	channelID   string
	players     []*Player
	started     bool
	turn        int
	lead        int
	current     Combo
	skipsInRow int
	winnerIdx   int
	gameOver    bool
	handSize    int
}

// New creates a new Game for a given channelID.
func New(channelID string) *Game {
	return &Game{
		channelID: channelID,
		players:   make([]*Player, 0, 4),
		winnerIdx: -1,
	}
}

// AddPlayer adds a player to the game.
func (g *Game) AddPlayer(userID, name, tag string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.started {
		return errors.New("game already started")
	}
	if len(g.players) >= 4 {
		return errors.New("game is full (max 4 players)")
	}
	for _, player := range g.players {
		if player.UserID == userID {
			return errors.New("you already joined")
		}
	}
	g.players = append(g.players, &Player{
		Idx:    len(g.players),
		UserID: userID,
		Name:   name,
		Tag:    tag,
	})
	return nil
}

// RemovePlayer removes player by userID.
func (g *Game) RemovePlayer(userID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for i, player := range g.players {
		if player.UserID == userID {
			g.players = append(g.players[:i], g.players[i+1:]...)
			for j := range g.players {
				g.players[j].Idx = j
			}
			break
		}
	}
}

// PlayerList returns a list of player names.
func (g *Game) PlayerList() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	names := make([]string, 0, len(g.players))
	for _, player := range g.players {
		names = append(names, player.Name)
	}
	return strings.Join(names, ", ")
}

// PlayersSnapshot returns a copy of the players slice.
func (g *Game) PlayersSnapshot() []*Player {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make([]*Player, len(g.players))
	copy(out, g.players)
	return out
}

// NumPlayers returns the number of players currently in the game.
func (g *Game) NumPlayers() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return len(g.players)
}

// FindPlayer returns a player's details by userID.
func (g *Game) FindPlayer(userID string) *Player {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, p := range g.players {
		if p.UserID == userID {
			return p
		}
	}
	return nil
}

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
		player.Skiped = false
		player.Finished = false
	}
	g.handSize = len(hands[0])

	start := 0
	found3S := false
	for i, p := range g.players {
		had3S := false
		p.Hand, had3S = discardThrees(p.Hand)
		if had3S {
			start = i
			found3S = true
		}
		SortCardsAsc(p.Hand)
	}
	if !found3S {
		start = 0
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

	cp := g.currentPlayer()
	if cp.UserID != userID {
		return "", errors.New("not your turn")
	}
	if len(codes) == 0 {
		return "", errors.New("provide cards to play, e.g.`play 2H 2S`")
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

	cp.RemoveCards(combo.Cards)
	g.current = combo
	g.lead = g.turn
	g.skipsInRow = 0

	if len(cp.Hand) == 0 {
		cp.Finished = true
		g.gameOver = true
		g.winnerIdx = cp.Idx
		return fmt.Sprintf("%s plays %s and wins the game.", cp.Name, comboString(combo)), nil
	}

	g.advanceTurn()
	return fmt.Sprintf("%s plays %s. Next: %s", cp.Name, comboString(combo), g.currentPlayer().Name), nil
}

// Skip handles a player skipping their turn.
func (g *Game) Skip(userID string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if !g.started {
		return "", errors.New("game not started")
	}
	if g.current.Type == ComboNone {
		return "", errors.New("cannot skip on an empty table")
	}

	cp := g.currentPlayer()
	if cp.UserID != userID {
		return "", errors.New("not your turn")
	}

	cp.Skiped = true
	g.skipsInRow++
	g.advanceTurn()

	if g.skipsInRow >= g.activePlayers()-1 {
		g.current = Combo{Type: ComboNone}
		g.skipsInRow = 0
		g.turn = g.lead
		return fmt.Sprintf("%s skips. Table cleared. %s to lead.", cp.Name, g.currentPlayer().Name), nil
	}
	return fmt.Sprintf("%s skips. Next: %s", cp.Name, g.currentPlayer().Name), nil
}

// TableStateString returns the current table state.
func (g *Game) TableStateString() string {
	if !g.started {
		return "Game not started."
	}
	if g.current.Type == ComboNone {
		return fmt.Sprintf("Table is empty. %s to lead.", g.currentPlayer().Name)
	}
	return fmt.Sprintf("On table: %s. %s to act.", comboString(g.current), g.currentPlayer().Name)
}

// IsOver returns whether the game is over.
func (g *Game) IsOver() bool {
	return g.gameOver
}

// WinnerName returns the name of the winner.
func (g *Game) WinnerName() string {
	if g.winnerIdx < 0 || g.winnerIdx >= len(g.players) {
		return ""
	}
	return g.players[g.winnerIdx].Name
}

func (g *Game) advanceTurn() {
	for {
		g.turn = (g.turn + 1) % len(g.players)
		if !g.players[g.turn].Finished {
			return
		}
	}
}

func (g *Game) activePlayers() int {
	n := 0
	for _, p := range g.players {
		if !p.Finished {
			n++
		}
	}
	return n
}

func (g *Game) currentPlayer() *Player {
	return g.players[g.turn]
}
