package store

import (
	"sync"

	"github.com/yoru0/odibot/internal/game"
)

// Session represents one lobby/game bound to a single guild text channel.
type Session struct {
	mu             sync.RWMutex
	Game           *game.Game
	Desired        int
	LobbyChannelID string
	DMChannel      map[string]string
	Selected       map[string][]string
	Started        bool
	HasDummy       bool
}

// newSession creates a Session with initialized maps.
func newSession(channelID string, desired int) *Session {
	return &Session{
		Game:           game.New(channelID),
		Desired:        desired,
		LobbyChannelID: channelID,
		DMChannel:      make(map[string]string, 4),
		Selected:       make(map[string][]string, 4),
	}
}

// GetDMChannel returns the DM channel ID for userID.
func (s *Session) GetDMChannel(userID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.DMChannel[userID]
}

// SetDMChannel sets the DM channel ID for userID.
func (s *Session) SetDMChannel(userID, channelID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DMChannel[userID] = channelID
}

// GetSelected returns a copy of the selected cards for userID.
func (s *Session) GetSelected(userID string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.Selected[userID]
	cp := make([]string, len(src))
	copy(cp, src)
	return cp
}

// SetSelected overwrites the selection for userID with a copy.
func (s *Session) SetSelected(userID string, selected []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := make([]string, len(selected))
	copy(cp, selected)
	s.Selected[userID] = cp
}

// DeleteSelected removes any selection for userID.
func (s *Session) DeleteSelected(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Selected, userID)
}
