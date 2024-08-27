package memstore

import (
	"context"
	"encoding/json"
	sessions "github.com/gowriter/sessions"
)

type MemoryStore struct {
	Store map[string]any
}

func NewMemoryStore() sessions.Store {
	return &MemoryStore{Store: map[string]any{}}
}

func (m *MemoryStore) New(_ context.Context, sessionID string, data any) error {
	m.Store[sessionID] = data
	return nil
}

func (m *MemoryStore) Get(_ context.Context, sessionID string) (json.RawMessage, error) {
	data, ok := m.Store[sessionID]
	if !ok {
		return nil, sessions.ErrNotFound
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (m *MemoryStore) Put(_ context.Context, sessionID string, session any) error {
	m.Store[sessionID] = session
	return nil
}

func (m *MemoryStore) End(_ context.Context, sessionID string) error {
	delete(m.Store, sessionID)
	return nil
}
