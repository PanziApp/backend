package usecase

import (
	"context"
	"github.com/PanziApp/backend/internal/domain"
)

type (
	UserRepository interface {
		Create(ctx context.Context, user domain.User) (userId domain.EntityId, err error)

		Get(ctx context.Context, userId domain.EntityId) (domain.User, error)
		GetByEmail(ctx context.Context, email domain.Email) (domain.User, error)

		Update(ctx context.Context, userId domain.EntityId, updates domain.EntityUpdate) error
	}

	SessionRepository interface {
		Create(ctx context.Context, session domain.Session) (sessionId domain.EntityId, err error)

		GetByToken(ctx context.Context, token domain.Token) (domain.Session, error)

		Update(ctx context.Context, sessionId domain.EntityId, updates domain.EntityUpdate) error
	}
)

type (
	Mailer interface {
		Send(ctx context.Context, receiver, name, subject, messageInHtml string) error
	}
)
