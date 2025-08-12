package game

import "strings"

type Player struct {
	Idx      int
	UserID   string
	Name     string
	Tag      string
	Hand     []Card
	Skiped   bool
	Finished bool
}

// HasCard returns card index, returns -1 if none is found.
func (p *Player) HasCard(c Card) int {
	for i, h := range p.Hand {
		if c == h {
			return i
		}
	}
	return -1
}

// RemoveCards removes list of selected cards.
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
