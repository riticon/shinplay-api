package session

import (
	"context"
	"time"

	"github.com/shinplay/ent"
	"github.com/shinplay/ent/session"
)

type SessionRepositoryIntr interface {
	CreateNewSession(ctx context.Context, user *ent.User, refreshToken string, expiresAt time.Time, userAgent string, ipAddress string) (*ent.Session, error)
	FindSessionByID(ctx context.Context, sessionID string) (*ent.Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
}

type SessionRepository struct {
	client *ent.Client
}

func NewSessionRepository(client *ent.Client) *SessionRepository {
	return &SessionRepository{client: client}
}

func (s *SessionRepository) CreateNewSession(ctx context.Context, user *ent.User, refreshToken string, expiresAt time.Time, userAgent, ipAddress string) (*ent.Session, error) {
	return s.client.Session.Create().
		SetUser(user).
		SetRefreshToken(refreshToken).
		SetExpiresAt(expiresAt).
		SetUserAgent(userAgent).
		SetIPAddress(ipAddress).
		Save(ctx)
}

func (s *SessionRepository) FindSessionByID(ctx context.Context, sessionID string) (*ent.Session, error) {
	return s.client.Session.Query().
		Where(session.SessionID(sessionID)).
		First(ctx)
}

func (s *SessionRepository) DeleteSession(ctx context.Context, sessionID string) (int, error) {
	id, err := s.client.Session.Delete().
		Where(session.SessionID(sessionID)).
		Exec(ctx)
	return id, err
}
