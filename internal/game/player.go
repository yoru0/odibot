package game

import "strings"

type Player struct {
	Idx      int
	UserID   string
	Name     string
	Tag      string
	Hand     []Card
	Skipped  bool
	Finished bool
	IsDummy  bool
}

// HasCard returns the index of card c, or -1 if not found.
func (p *Player) HasCard(c Card) int {
	for i, h := range p.Hand {
		if c == h {
			return i
		}
	}
	return -1
}

// RemoveCards removes the listed cards from the player's hand.
func (p *Player) RemoveCards(selected []Card) {
	for _, c := range selected {
		for i := range p.Hand {
			if p.Hand[i] == c {
				p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
				break
			}
		}
	}
}

// HandString returns a player's hand.
func (p *Player) HandString() string {
	var s strings.Builder
	for i, c := range p.Hand {
		if i > 0 && i < len(p.Hand) {
			s.WriteString(" | ")
		}
		s.WriteString(c.String())
	}
	return s.String()
}
