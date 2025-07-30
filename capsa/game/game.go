package game

import "github.com/yoru0/odibot/capsa/player"

type Game struct {
	Players       player.Players
	CurrIndex     int
	Round         int
	PlayerSkipped int
}
