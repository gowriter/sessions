package session_test

import (
	"context"
	"github.com/gowriter/session"
	"github.com/gowriter/session/memstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

const testSessionID = "test"

type testType struct {
	Id   string `mapstructure:"id" json:"id"`
	Name string `mapstructure:"name" json:"name"`
}

func TestNewSession(t *testing.T) {
	testStore := memstore.NewMemoryStore()
	ctx := context.Background()
	testTime := time.Now().Add(10 * time.Second)

	testCookie := http.Cookie{
		Name:     "test",
		Value:    "testID",
		Expires:  testTime,
		MaxAge:   100,
		HttpOnly: true,
	}

	testData := testType{
		Id:   "some",
		Name: "name",
	}

	t.Run("success new session", func(t *testing.T) {
		httpSession, err := session.NewHttpSession[testType](testStore, testCookie, session.Options{
			GeneratorFunc: func() (string, error) {
				return testSessionID, nil
			},
			RegenerateIDs: false,
		})
		require.NoError(t, err)

		returnedSession, err := httpSession.NewSession(ctx, testData)
		require.NoError(t, err)

		expectedCookie := http.Cookie{
			Name:     "test",
			Value:    testSessionID, // set in generator func
			Expires:  testTime,
			MaxAge:   100,
			HttpOnly: true,
		}

		assert.Equal(t, testData, *returnedSession.Object)
		assert.Equal(t, expectedCookie, *returnedSession.Cookie)
	})

	t.Run("success put session", func(t *testing.T) {
		httpSession, err := session.NewHttpSession[testType](testStore, testCookie, session.Options{
			GeneratorFunc: func() (string, error) {
				return "testID", nil
			},
			RegenerateIDs: false,
		})
		require.NoError(t, err)

		testSession := &session.Session[testType]{
			Object: &testData,
			Cookie: &testCookie,
		}

		err = httpSession.PutSession(ctx, testSession)
		assert.NoError(t, err)
	})

	t.Run("success put session", func(t *testing.T) {
		httpSession, err := session.NewHttpSession[testType](testStore, testCookie, session.Options{
			GeneratorFunc: func() (string, error) {
				return "testID", nil
			},
			RegenerateIDs: false,
		})
		require.NoError(t, err)

		wantSession := &session.Session[testType]{
			Object: &testData,
			Cookie: &testCookie,
		}

		gotSession, err := httpSession.GetSession(ctx, wantSession)
		assert.NoError(t, err)

		assert.Equal(t, wantSession, gotSession)
	})

	t.Run("success end session", func(t *testing.T) {
		httpSession, err := session.NewHttpSession[testType](testStore, testCookie, session.Options{
			GeneratorFunc: func() (string, error) {
				return "testID", nil
			},
			RegenerateIDs: false,
		})
		require.NoError(t, err)

		testSession := &session.Session[testType]{
			Object: &testData,
			Cookie: &testCookie,
		}

		gotSession, err := httpSession.EndSession(ctx, testSession)
		assert.NoError(t, err)

		_ = gotSession
	})
}
