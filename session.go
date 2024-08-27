package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type HttpSession[T any] struct {
	lock       sync.Mutex
	store      Store
	cookie     http.Cookie
	generate   func() (string, error)
	regenerate bool
}

type Session[T any] struct {
	Object *T
	Cookie *http.Cookie
}

type GeneratorFunc func() (string, error)

type Options struct {
	GeneratorFunc func() (string, error)
	RegenerateIDs bool
}

func NewHttpSession[T any](store Store, cookie http.Cookie, options Options) (*HttpSession[T], error) {
	httpSession := &HttpSession[T]{
		lock:       sync.Mutex{},
		store:      store,
		cookie:     cookie,
		generate:   options.GeneratorFunc,
		regenerate: options.RegenerateIDs,
	}

	return httpSession, nil
}

func (service *HttpSession[T]) NewSession(ctx context.Context, t T) (*Session[T], error) {
	var sessionID string
	if service.generate != nil {
		id, err := service.generate()
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate session ID")
		}
		sessionID = id
	}

	err := service.store.New(ctx, sessionID, t)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store new session")
	}

	cookie := service.cookie
	cookie.Value = sessionID

	return &Session[T]{
		Object: &t,
		Cookie: &cookie,
	}, nil
}

func (service *HttpSession[T]) GetSession(ctx context.Context, session *Session[T]) (*Session[T], error) {
	service.lock.Lock()
	defer service.lock.Unlock()

	sessionData, err := service.store.Get(ctx, session.Cookie.Value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Data from session")
	}

	obj := new(T)
	err = json.Unmarshal(sessionData, obj)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal session Data")
	}

	session.Object = obj
	return session, nil
}

func (service *HttpSession[T]) PutSession(ctx context.Context, session *Session[T]) error {
	service.lock.Lock()
	defer service.lock.Unlock()

	err := service.store.Put(ctx, session.Cookie.Value, session.Object)
	if err != nil {
		return errors.Wrap(err, "failed to put session in storage")
	}

	return nil
}

func (service *HttpSession[T]) EndSession(ctx context.Context, session *Session[T]) (*Session[T], error) {
	service.lock.Lock()
	defer service.lock.Unlock()

	err := service.store.End(ctx, session.Cookie.Value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to delete session")
	}

	session.Cookie.MaxAge = -1
	session.Cookie.Expires = time.Now()
	session.Cookie.Value = ""
	session.Object = nil
	return session, nil
}

func WithGenerateUUID() GeneratorFunc {
	return func() (string, error) {
		return uuid.New().String(), nil
	}
}

func WithGenerateRBID() GeneratorFunc {
	return func() (string, error) {
		b := make([]byte, 16)
		_, err := io.ReadFull(rand.Reader, b)
		if err != nil {
			return "", err
		}
		return base64.URLEncoding.EncodeToString(b), nil
	}
}
