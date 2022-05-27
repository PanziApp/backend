package repo

import (
	"context"
	"github.com/PanziApp/backend/internal/domain"
	"github.com/PanziApp/backend/pkg/postgres"
)

type UserRepository struct {
	postgres.Postgres
}

func NewUserRepository(pg postgres.Postgres) UserRepository {
	return UserRepository{pg}
}

func (r UserRepository) Create(ctx context.Context, u domain.User) (domain.EntityId, error) {
	sql, args, err := r.Builder.
		Insert("users").
		Columns("create_time, email, email_verify_time, hashed_password, fullname, avatar").
		Values(u.CreateTime, u.Email, u.EmailVerifyTime, u.HashedPassword, u.Fullname, u.Avatar).
		Suffix("returning id").
		ToSql()
	if err != nil {
		return 0, domain.InternalError{Err: err}
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&u.Id)
	if err != nil {
		return 0, domain.InternalError{Err: err}
	}
	return u.Id, nil
}

func (r UserRepository) Get(ctx context.Context, userId domain.EntityId) (u domain.User, err error) {
	sql, args, err := r.Builder.
		Select("id, create_time, email, email_verify_time, hashed_password, fullname, avatar").
		From("users").
		Where("id = ?", userId).
		ToSql()
	if err != nil {
		return u, domain.InternalError{Err: err}
	}

	err = r.Pool.QueryRow(ctx, sql, args...).
		Scan(&u.Id, &u.CreateTime, &u.Email, &u.EmailVerifyTime, &u.HashedPassword, &u.Fullname, &u.Avatar)
	if err != nil {
		return u, domain.InternalError{Err: err}
	}
	return u, nil
}

func (r UserRepository) GetByEmail(ctx context.Context, email domain.Email) (u domain.User, err error) {
	sql, args, err := r.Builder.
		Select("id, create_time, email, email_verify_time, hashed_password, fullname, avatar").
		From("users").
		Where("email = ?", email).
		ToSql()
	if err != nil {
		return u, domain.InternalError{Err: err}
	}

	err = r.Pool.QueryRow(ctx, sql, args...).
		Scan(&u.Id, &u.CreateTime, &u.Email, &u.EmailVerifyTime, &u.HashedPassword, &u.Fullname, &u.Avatar)
	if err != nil {
		return u, domain.InternalError{Err: err}
	}
	return u, nil
}

func (r UserRepository) Update(ctx context.Context, userId domain.EntityId, updates domain.EntityUpdate) error {
	q := r.Builder.Update("users").
		Where("id = ?", userId)

	haveUpdate := false
	if emailVerifyTime, ok := updates[domain.UserEmailVerifyTimeFieldName]; ok {
		q.Set("email_verify_time", emailVerifyTime)
		haveUpdate = true
	}
	if hashedPassword, ok := updates[domain.UserHashedPasswordFieldName]; ok {
		q.Set("hashed_password", hashedPassword)
		haveUpdate = true
	}
	if fullname, ok := updates[domain.UserFullnameFieldName]; ok {
		q.Set("fullname", fullname)
		haveUpdate = true
	}
	if avatar, ok := updates[domain.UserAvatarFieldName]; ok {
		q.Set("avatar", avatar)
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
