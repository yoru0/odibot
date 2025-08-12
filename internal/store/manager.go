package store

import (
	"sync"

	"github.com/yoru0/odibot/internal/game"
)

type Manager struct {
	mu        sync.RWMutex
	byChannel map[string]*Session
	byUser    map[string]*Session
}

func NewManager() *Manager {
	return &Manager{
		byChannel: make(map[string]*Session),
		byUser:    make(map[string]*Session),
	}
}

func (m *Manager) NewSession(channelID string, desired int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.byChannel[channelID] = &Session{
		LobbyChannelID: channelID,
		Desired:        desired,
		Game:           game.New(channelID),
		DMChannel:      make(map[string]string),
	}
}

// Get returns the details of a Session by channelID.
func (m *Manager) Get(channelID string) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.byChannel[channelID]
}

// GetByUser returns the details of a Session by userID.
func (m *Manager) GetByUser(userID string) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.byUser[userID]
}

// Has returns whether a Session is in a channelID.
func (m *Manager) Has(channelID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.byChannel[channelID]
	return ok
}

// MarkStarted marks started Channel and index all players.
func (m *Manager) MarkStarted(channelID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s := m.byChannel[channelID]
	if s == nil {
		return
	}
	s.Started = true
	for _, p := range s.Game.PlayersSnapshot() {
		m.byUser[p.UserID] = s
	}
}

// Delete deletes a Session by channelID.
func (m *Manager) Delete(channelID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s := m.byChannel[channelID]
	if s != nil {
		for _, p := range s.Game.PlayersSnapshot() {
			delete(m.byUser, p.UserID)
		}
		delete(m.byChannel, channelID)
	}
}
