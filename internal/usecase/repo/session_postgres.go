package repo

import (
	"context"
	"github.com/PanziApp/backend/internal/domain"
	"github.com/PanziApp/backend/pkg/postgres"
)

type SessionRepository struct {
	postgres.Postgres
}

func NewSessionRepository(pg postgres.Postgres) SessionRepository {
	return SessionRepository{pg}
}

func (r SessionRepository) Create(ctx context.Context, s domain.Session) (domain.EntityId, error) {
	sql, args, err := r.Builder.
		Insert("sessions").
		Columns("create_time, user_id, type, token, valid_until").
		Values(s.CreateTime, s.UserId, s.Type, s.Token, s.ValidUntil).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, domain.InternalError{Err: err}
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&s.Id)
	if err != nil {
		return 0, domain.InternalError{Err: err}
	}
	return s.Id, nil
}

func (r SessionRepository) GetByToken(ctx context.Context, token domain.Token) (s domain.Session, err error) {
	sql, args, err := r.Builder.
		Select("id, create_time, user_id, type, token, valid_until").
		From("sessions").
		Where("token = ?", token).
		ToSql()
	if err != nil {
		return s, domain.InternalError{Err: err}
	}

	err = r.Pool.QueryRow(ctx, sql, args...).
		Scan(&s.Id, &s.CreateTime, &s.UserId, &s.Type, &s.Type, &s.Token, &s.ValidUntil)
	if err != nil {
		return s, domain.InternalError{Err: err}
	}
	return s, nil
}

func (r SessionRepository) Update(ctx context.Context, sessionId domain.EntityId, updates domain.EntityUpdate) error {
	q := r.Builder.Update("sessions").
		Where("id = ?", sessionId)

	haveUpdate := false
	if validUntil, ok := updates[domain.SessionValidUntilFieldName]; ok {
		q.Set("valid_until", validUntil)
		haveUpdate = true
	}

	if !haveUpdate {
		return nil
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return domain.InternalError{Err: err}
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return domain.InternalError{Err: err}
	}

	return nil
}
