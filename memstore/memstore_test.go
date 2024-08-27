package memstore

import (
	"context"
	"encoding/json"
	sessions "github.com/gowriter/sessions"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestType struct {
	ID   int
	Name string
}

func TestMemoryStore_Get(t *testing.T) {
	type fields struct {
		Store sessions.Store
	}
	type args struct {
		ctx       context.Context
		sessionID string
	}

	ctx := context.Background()
	testSessionData := TestType{
		ID:   1,
		Name: "test",
	}

	testSessionJSON, err := json.Marshal(testSessionData)
	require.NoError(t, err)

	tests := []struct {
		name        string
		sessionJSON []byte
		fields      fields
		args        args
		want        json.RawMessage
		wantErr     error
	}{
		{
			name: "sessions exists",
			fields: fields{
				Store: NewMemoryStore(),
			},
			args: args{
				ctx:       ctx,
				sessionID: "test_session_id",
			},
			want: testSessionJSON,
		},
		{
			name: "error sessions not exists",
			fields: fields{
				Store: NewMemoryStore(),
			},
			args: args{
				ctx:       ctx,
				sessionID: "wrong_test_session_id",
			},
			wantErr: sessions.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.Store.(*MemoryStore)
			s.Store["test_session_id"] = testSessionData

			got, err := tt.fields.Store.Get(tt.args.ctx, tt.args.sessionID)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStore_Put(t *testing.T) {
	type fields struct {
		Store sessions.Store
	}
	type args struct {
		ctx         context.Context
		sessionID   string
		sessionData TestType
	}

	ctx := context.Background()

	tests := []struct {
		name        string
		sessionJSON []byte
		fields      fields
		args        args
		wantErr     bool
	}{
		{
			name: "success put value",
			fields: fields{
				Store: NewMemoryStore(),
			},
			args: args{
				ctx:       ctx,
				sessionID: "test_session_id",
				sessionData: TestType{
					ID:   1,
					Name: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.Store.(*MemoryStore)
			s.Store[tt.args.sessionID] = TestType{}

			err := s.Put(tt.args.ctx, tt.args.sessionID, tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.EqualError(t, err, sessions.ErrNotFound.Error())
			}
			assert.NoError(t, err)
			got, ok := s.Store[tt.args.sessionID]
			assert.True(t, ok)
			assert.Equal(t, tt.args.sessionData, got)
		})
	}
}

func TestMemoryStore_End(t *testing.T) {
	type fields struct {
		Store sessions.Store
	}
	type args struct {
		ctx       context.Context
		sessionID string
	}

	ctx := context.Background()

	tests := []struct {
		name        string
		sessionJSON []byte
		fields      fields
		args        args
		wantErr     bool
	}{
		{
			name: "success end sessions",
			fields: fields{
				Store: NewMemoryStore(),
			},
			args: args{
				ctx:       ctx,
				sessionID: "test_session_id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.Store.(*MemoryStore)
			s.Store[tt.args.sessionID] = TestType{}

			err := s.End(tt.args.ctx, tt.args.sessionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("End() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.EqualError(t, err, sessions.ErrNotFound.Error())
			}
			assert.NoError(t, err)
			_, ok := s.Store[tt.args.sessionID]
			assert.False(t, ok)
		})
	}
}

func TestMemoryStore_New(t *testing.T) {
	type fields struct {
		Store sessions.Store
	}
	type args struct {
		ctx         context.Context
		sessionID   string
		sessionData TestType
	}

	ctx := context.Background()

	tests := []struct {
		name        string
		sessionJSON []byte
		fields      fields
		args        args
		wantErr     bool
	}{
		{
			name: "success new sessions",
			fields: fields{
				Store: NewMemoryStore(),
			},
			args: args{
				ctx:       ctx,
				sessionID: "test_session_id",
				sessionData: TestType{
					ID:   1,
					Name: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.Store.(*MemoryStore)

			err := s.New(tt.args.ctx, tt.args.sessionID, tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.EqualError(t, err, sessions.ErrNotFound.Error())
			}
			assert.NoError(t, err)
			got, ok := s.Store[tt.args.sessionID]
			assert.True(t, ok)
			assert.Equal(t, tt.args.sessionData, got)
		})
	}
}
