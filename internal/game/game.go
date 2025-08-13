package game

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Game struct {
	mu            sync.Mutex
	channelID     string
	players       []*Player
	started       bool
	turn          int
	lead          int
	current       Combo
	skipsInRow    int
	winnerIdx     int
	gameOver      bool
	handSize      int
	finishedOrder []int
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

// AddDummy adds a dummy to the game (owner-controlled).
func (g *Game) AddDummy(userID, name string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.started {
		return errors.New("game already started")
	}
	if len(g.players) >= 4 {
		return errors.New("game is full (max 4 players)")
	}
	for _, p := range g.players {
		if p.UserID == userID {
			return errors.New("user already joined")
		}
	}
	g.players = append(g.players, &Player{
		Idx:     len(g.players),
		UserID:  userID,
		Name:    name,
		Tag:     "dummy",
		IsDummy: true,
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

// CurrentPlayerInfo returns information about current player.
func (g *Game) CurrentPlayerInfo() (id, name string, isDummy bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	p := g.players[g.turn]
	return p.UserID, p.Name, p.IsDummy
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

// ResultsString returns winner's order.
func (g *Game) ResultsString() string {
	if len(g.players) == 0 {
		return ""
	}
	order := make([]string, 0, len(g.players))
	for _, idx := range g.finishedOrder {
		order = append(order, g.players[idx].Name)
	}
	for _, p := range g.players {
		if !p.Finished {
			order = append(order, p.Name)
			break
		}
	}
	var b strings.Builder
	for i, name := range order {
		fmt.Fprintf(&b, "%d. %s\n", i+1, name)
	}
	return strings.TrimRight(b.String(), "\n")
}

// WinnerName returns the name of the winner.
func (g *Game) WinnerName() string {
	if len(g.finishedOrder) > 0 {
		i := g.finishedOrder[0]
		if i >= 0 && i < len(g.players) {
			return g.players[i].Name
		}
	}
	if g.winnerIdx >= 0 && g.winnerIdx < len(g.players) {
		return g.players[g.winnerIdx].Name
	}
	return ""
}

func (g *Game) FormatThreesReport() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	spadeIdx := -1
	for i, p := range g.players {
		if p.HasCard(Card{Rank: R3, Suit: Spades}) >= 0 {
			spadeIdx = i
			break
		}
	}
	var b strings.Builder
	b.WriteString("```")
	for i, p := range g.players {
		prefix := "   "
		if i == spadeIdx {
			prefix = "-> "
		}

		parts := make([]string, 0, 4)
		if p.HasCard(Card{Rank: R3, Suit: Diamonds}) >= 0 {
			parts = append(parts, "3D")
		}
		if p.HasCard(Card{Rank: R3, Suit: Clubs}) >= 0 {
			parts = append(parts, "3C")
		}
		if p.HasCard(Card{Rank: R3, Suit: Hearts}) >= 0 {
			parts = append(parts, "3H")
		}
		if p.HasCard(Card{Rank: R3, Suit: Spades}) >= 0 {
			parts = append(parts, "3S")
		}
		if len(parts) != 0 {
			b.WriteString(fmt.Sprintf("%s%s - %s\n", prefix, p.Name, strings.Join(parts, " ")))
		}
		p.Hand = discardThrees(p.Hand)
	}
	if spadeIdx >= 0 {
		b.WriteString(fmt.Sprintf("\n%s plays first", g.players[spadeIdx].Name))
	}
	b.WriteString("```")
	return b.String()
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
