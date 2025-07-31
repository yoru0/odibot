package player

type SkipedPlayerName []string

func (s *SkipedPlayerName) ResetSkippedPlayer() {
	*s = []string{}
}
