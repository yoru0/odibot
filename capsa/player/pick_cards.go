package player

import "sort"

func (p *Player) RemovePlayedCards(picks []int) {
	sort.Sort(sort.Reverse(sort.IntSlice(picks)))
	for _, i := range picks {
		if i >= 0 && i < len(p.Hand) {
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
		}
	}
}
