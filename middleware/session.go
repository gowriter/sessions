package middleware

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

type sessionDataType int

var SessionDataKey sessionDataType

func NewSessionMiddleware[T any](h http.Handler, l logrus.FieldLogger, httpSession *session.HttpSession[T], cName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := &session.Session[T]{}
		c, err := r.Cookie(cName)
		switch {
		case err != nil:
			var genericT T
			s, err = httpSession.NewSession(r.Context(), genericT)
			if err != nil {
				http.Error(w, "failed to create new session", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, c)
		default:
			s, err = httpSession.GetSession(r.Context(), s)
			if err != nil {
				http.Error(w, "failed to get session", http.StatusInternalServerError)
				return
			}
		}

		r = r.WithContext(ToContext(r.Context(), s.Object))
		h.ServeHTTP(w, r)

		if err = httpSession.PutSession(r.Context(), s); err != nil {
			l.WithError(err).Error("failed to put session data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func FromContext[T any](ctx context.Context) *T {
	val := ctx.Value(SessionDataKey)
	t, ok := val.(*T)
	if !ok {
		return nil
	}

	return t
}

func ToContext(ctx context.Context, t any) context.Context {
	return context.WithValue(ctx, SessionDataKey, t)
}
