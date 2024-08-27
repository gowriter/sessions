package session

import (
	"context"
	"encoding/json"
	"errors"
)

var ErrNotFound = errors.New("not found")

type Store interface {
	New(ctx context.Context, sessionID string, Object any) error
	Get(ctx context.Context, sessionID string) (json.RawMessage, error)
	Put(ctx context.Context, sessionID string, Object any) error
	End(ctx context.Context, sessionID string) error
}
