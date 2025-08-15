package store

import (
	"sync"
)

// Manager indexes all active sessions by lobby channel and by user ID.
type Manager struct {
	mu        sync.RWMutex
	byChannel map[string]*Session
	byUser    map[string]*Session
}

// NewManager creates an empty Manager.
func NewManager() *Manager {
	return &Manager{
		byChannel: make(map[string]*Session),
		byUser:    make(map[string]*Session),
	}
}

// NewSession creates a session for channelID.
func (m *Manager) NewSession(channelID string, desired int) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.byChannel[channelID]; ok {
		return s
	}
	sess := newSession(channelID, desired)
	m.byChannel[channelID] = sess
	return sess
}

// Has returns whether a lobby exists for the given channel.
func (m *Manager) Has(channelID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.byChannel[channelID]
	return ok
}

// Get returns the Session for a guild channelID.
func (m *Manager) Get(channelID string) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.byChannel[channelID]
}

// GetByUser returns the started Session a user belongs to.
func (m *Manager) GetByUser(userID string) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.byUser[userID]
}

// MarkStarted marks the lobby as started and indexes all players by userID.
func (m *Manager) MarkStarted(channelID, ownerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := m.byChannel[channelID]
	if s == nil {
		return
	}
	s.Started = true

	players := s.Game.PlayersSnapshot()
	for _, p := range players {
		m.byUser[p.UserID] = s
	}

	// for _, p := range s.Game.PlayersSnapshot() {
	// 	m.byUser[p.UserID] = s
	// }
	// if s.HasDummy && ownerID != "" {
	// 	m.byUser[ownerID] = s
	// }
}

// Delete removes a Session by its lobby channelID.
func (m *Manager) Delete(channelID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := m.byChannel[channelID]
	if s == nil {
		return
	}

	for _, p := range s.Game.PlayersSnapshot() {
		if m.byUser[p.UserID] == s {
			delete(m.byUser, p.UserID)
		}
	}
	delete(m.byChannel, channelID)
}
