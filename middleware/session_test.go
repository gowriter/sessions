package middleware_test

import (
	"context"
	"encoding/json"
	"github.com/gowriter/sessions/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gowriter/sessions/memstore"
)

// json tags has to be set for correct encoding and decoding
type TestSessionType struct {
	Id   string `mapstructure:"id" json:"id"`
	Name string `mapstructure:"name" json:"name"`
}

func TestNewSessionHandler(t *testing.T) {
	const testSessionID = "testID"
	var store session.Store

	ctx := context.Background()
	l, _ := test.NewNullLogger()

	inputData := TestSessionType{
		Id:   "some_id",
		Name: "some_name",
	}

	outputData := TestSessionType{
		Id:   "new_id",
		Name: "new_name",
	}

	t.Run("success valid session", func(t *testing.T) {
		store = memstore.NewMemoryStore()
		testSession := &session.Session[TestSessionType]{
			Object: &TestSessionType{
				Id:   "1",
				Name: "test",
			},
			Cookie: &http.Cookie{
				Name:     "test",
				Value:    testSessionID,
				Expires:  time.Now().Add(10 * time.Second),
				MaxAge:   100,
				HttpOnly: true,
			},
		}

		router := http.NewServeMux()
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			sess := middleware.FromContext[TestSessionType](r.Context())
			assert.Equal(t, &inputData, sess)

			// set new values to store
			sess.Id = outputData.Id
			sess.Name = outputData.Name
		})

		httpSession, err := session.NewHttpSession[TestSessionType](store, *testSession.Cookie, session.Options{
			GeneratorFunc: func() (string, error) {
				return testSessionID, nil
			},
		})
		require.NoError(t, err)
		sessionRouter := middleware.NewSessionMiddleware[TestSessionType](router, logrus.NewEntry(l), httpSession, "test")

		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", http.NoBody).WithContext(ctx)

		err = store.Put(ctx, testSessionID, testSession)
		require.NoError(t, err)

		req.AddCookie(testSession.Cookie)
		sessionRouter.ServeHTTP(rw, req)

		got, err := store.Get(ctx, testSessionID)
		require.NoError(t, err)

		b, err := json.Marshal(outputData)
		require.NoError(t, err)

		output := json.RawMessage{}
		err = json.Unmarshal(b, &output)
		require.NoError(t, err)

		assert.Equal(t, output, got)
	})

	//t.Run("no session cookie set", func(t *testing.T) {
	//	store = memstore.NewMemoryStore()
	//	router := http.NewServeMux()
	//	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//		sess := middleware.FromContext[TestSessionType](r.Context())
	//		assert.Empty(t, sess)
	//	})
	//
	//	httpSession, err := session.NewHttpSession[TestSessionType](store, *testSession.Cookie)
	//	require.NoError(t, err)
	//	sessionRouter := middleware.NewSessionMiddleware[TestSessionType](router, logrus.NewEntry(l), httpSession, "test")
	//
	//	rw := httptest.NewRecorder()
	//	req := httptest.NewRequest("GET", "/", http.NoBody).WithContext(ctx)
	//	sessionRouter.ServeHTTP(rw, req)
	//})

	//t.Run("error mapstructure decode 500", func(t *testing.T) {
	//	store = memstore.NewMemoryStore()
	//	sess := &session.Session{
	//		Cookie: &http.Cookie{
	//			Name:     memstore.SessionCookie,
	//			Value:    testSessionID,
	//			Expires:  time.Now().Add(10 * time.Second),
	//			MaxAge:   100,
	//			HttpOnly: true,
	//		},
	//	}
	//
	//	m := `{"id": 3, name": test}`
	//
	//	var err error
	//	sess.Data, err = json.Marshal([]byte(m))
	//	require.NoError(t, err)
	//
	//	err = store.Put(ctx, testSessionID, sess)
	//	require.NoError(t, err)
	//
	//	router := http.NewServeMux()
	//	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//		sess := middleware.FromContext[TestSessionType](r.Context())
	//		spew.Dump(sess)
	//	})
	//
	//	sessionRouter := middleware.NewSessionHandler[TestSessionType](router, logrus.NewEntry(l), store, memstore.SessionCookie)
	//
	//	rw := httptest.NewRecorder()
	//	req := httptest.NewRequest("GET", "/", http.NoBody).WithContext(ctx)
	//
	//	req.AddCookie(sess.Cookie)
	//	sessionRouter.ServeHTTP(rw, req)
	//
	//	assert.Equal(t, rw.Code, http.StatusInternalServerError)
	//})
	//
	//t.Run("error json unmarshal 500", func(t *testing.T) {
	//	store = memstore.NewMemoryStore()
	//	sess := &session.Session{
	//		Cookie: &http.Cookie{
	//			Name:     memstore.SessionCookie,
	//			Value:    testSessionID,
	//			Expires:  time.Now().Add(10 * time.Second),
	//			MaxAge:   100,
	//			HttpOnly: true,
	//		},
	//	}
	//
	//	sess.Data = []byte(nil)
	//	err := store.Put(ctx, testSessionID, sess)
	//	require.NoError(t, err)
	//
	//	router := http.NewServeMux()
	//	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//		sess := middleware.FromContext[TestSessionType](r.Context())
	//		assert.Empty(t, sess)
	//	})
	//
	//	sessionRouter := middleware.NewSessionHandler[TestSessionType](router, logrus.NewEntry(l), store, memstore.SessionCookie)
	//
	//	rw := httptest.NewRecorder()
	//	req := httptest.NewRequest("GET", "/", http.NoBody).WithContext(ctx)
	//
	//	req.AddCookie(sess.Cookie)
	//	sessionRouter.ServeHTTP(rw, req)
	//
	//	assert.Equal(t, rw.Code, http.StatusInternalServerError)
	//})
}
