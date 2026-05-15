package sessions

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/stpotter16/hours/internal/store/db"
)

type contextKey struct {
	name string
}

const SESSION_KEY = "session"
const SESSION_ENV_KEY = "HOURS_SESSION_ENV_KEY"

type Session struct {
	ID     string
	UserId int
}

type SessionManger struct {
	db                   *db.DB
	sessionHmacSecretKey string
}

func New(db db.DB, getenv func(string) string) (SessionManger, error) {
	hmacSecret := getenv(SESSION_ENV_KEY)
	if hmacSecret == "" {
		return SessionManger{}, errors.New("could not locate HMAC secret key")
	}

	s := SessionManger{
		db:                   &db,
		sessionHmacSecretKey: hmacSecret,
	}

	go s.sessionCleanup()

	return s, nil
}

func (s SessionManger) CreateSession(w http.ResponseWriter, r *http.Request, userId int) error {
	session := Session{
		ID:     uuid.NewString(),
		UserId: userId,
	}

	if err := s.writeSessionCookie(w, session); err != nil {
		return err
	}

	serializedSession, err := serializeSession(session)
	if err != nil {
		return err
	}

	return s.insertSession(r.Context(), session.ID, serializedSession)
}

func (s SessionManger) DeleteSession(w http.ResponseWriter, r *http.Request) error {
	session, err := s.loadSession(r)
	if err != nil {
		return err
	}
	if err = s.deleteSession(r.Context(), session.ID); err != nil {
		return err
	}
	return s.deleteSessionCookie(w)
}

func (s SessionManger) PopulateSessionContext(r *http.Request) (context.Context, error) {
	session, err := s.loadSession(r)
	if err != nil {
		return nil, err
	}

	ctxKey := contextKey{SESSION_KEY}
	return context.WithValue(r.Context(), ctxKey, session), nil
}

func (s SessionManger) SessionFromContext(ctx context.Context) (Session, error) {
	ctxKey := contextKey{SESSION_KEY}
	session, okay := ctx.Value(ctxKey).(Session)
	if !okay {
		return Session{}, errors.New("no session info in context")
	}
	return session, nil
}

func (s SessionManger) loadSession(r *http.Request) (Session, error) {
	cookie, err := s.readSessionCookie(r)
	if err != nil {
		return Session{}, err
	}

	cookieVals := strings.SplitN(cookie, "::", 2)
	if len(cookieVals) != 2 {
		return Session{}, errors.New("cookie is invalid")
	}
	cookieToken := cookieVals[1]

	serializedSession, err := s.readSession(r.Context(), cookieToken)
	if err != nil {
		return Session{}, err
	}

	session, err := deserializeSession(serializedSession)
	if err != nil {
		return Session{}, err
	}

	if cookieToken != session.ID {
		return Session{}, errors.New("invalid session token")
	}

	return session, nil
}

