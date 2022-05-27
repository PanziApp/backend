package usecase

import (
	"context"
	"github.com/PanziApp/backend/internal/domain"
	"time"
)

type UserUseCase struct {
	repo struct {
		user    UserRepository
		session SessionRepository
	}
	mailer Mailer
}

func New(
	userRepository UserRepository,
	sessionRepository SessionRepository,
	mailer Mailer,
) UserUseCase {
	uc := UserUseCase{}

	uc.repo.user = userRepository
	uc.repo.session = sessionRepository

	uc.mailer = mailer

	return uc
}

func (uc UserUseCase) createSession(
	ctx context.Context,
	userId domain.EntityId,
	tokenType domain.TokenType,
	validUntil *time.Time,
) (session domain.Session, err error) {
	session = domain.Session{
		CreateTime: time.Now(),
		UserId:     userId,
		Type:       tokenType,
		ValidUntil: validUntil,
	}
	session.Token, err = domain.RandomToken()
	if err != nil {
		return
	}
	session.Id, err = uc.repo.session.Create(ctx, session)
	return
}

func (uc UserUseCase) getGeneralValidSession(
	ctx context.Context,
	token string,
) (s domain.Session, err error) {
	validToken, err := domain.ValidateToken(token)
	if err != nil {
		return s, err
	}

	s, err = uc.repo.session.GetByToken(ctx, validToken)
	if err != nil {
		return s, err
	}

	if s.Type != domain.GeneralToken || (s.ValidUntil != nil && s.ValidUntil.Before(time.Now())) {
		return s, domain.ErrInvalidToken
	}

	return s, nil
}

func (uc UserUseCase) SignUp(
	ctx context.Context,
	email, password string,
) (token string, err error) {
	validEmail, err := domain.ValidateEmail(email)
	if err != nil {
		return "", err
	}

	validPassword, err := domain.ValidatePassword(password)
	if err != nil {
		return "", err
	}

	user := domain.User{
		CreateTime: time.Now(),
		Email:      validEmail,
	}
	user.HashedPassword, err = domain.HashPassword(validPassword)
	if err != nil {
		return "", err
	}
	user.Id, err = uc.repo.user.Create(ctx, user)
	if err != nil {
		return "", err
	}

	session, err := uc.createSession(ctx, user.Id, domain.EmailVerificationToken, nil)
	if err != nil {
		return "", err
	}

	err = uc.mailer.Send(
		ctx,
		string(user.Email),
		"User",
		"Email Verification",
		domain.EmailVerificationMessage(string(session.Token)),
	)
	if err != nil {
		return "", err
	}

	session, err = uc.createSession(ctx, user.Id, domain.GeneralToken, nil)
	if err != nil {
		return "", err
	}

	return string(session.Token), nil
}

func (uc UserUseCase) SignIn(
	ctx context.Context,
	email, password string,
) (token string, err error) {
	validEmail, err := domain.ValidateEmail(email)
	if err != nil {
		return "", err
	}

	validPassword, err := domain.ValidatePassword(password)
	if err != nil {
		return "", err
	}

	user, err := uc.repo.user.GetByEmail(ctx, validEmail)
	if err != nil {
		return "", err
	}

	err = user.HashedPassword.Match(validPassword)
	if err != nil {
		return "", err
	}

	session, err := uc.createSession(ctx, user.Id, domain.GeneralToken, nil)
	if err != nil {
		return "", err
	}

	return string(session.Token), nil
}

func (uc UserUseCase) SendResetPasswordLink(
	ctx context.Context,
	email string,
) error {
	validEmail, err := domain.ValidateEmail(email)
	if err != nil {
		return nil
	}

	user, err := uc.repo.user.GetByEmail(ctx, validEmail)
	if err != nil {
		return err
	}

	anHour := time.Now().Add(time.Hour)
	session, err := uc.createSession(ctx, user.Id, domain.ResetPasswordToken, &anHour)
	if err != nil {
		return err
	}

	err = uc.mailer.Send(
		ctx,
		string(user.Email),
		string(user.Fullname),
		"Reset Password",
		domain.ResetPasswordEmailMessage(string(session.Token)),
	)
	if err != nil {
		return err
	}

	return nil
}

func (uc UserUseCase) ResetPassword(
	ctx context.Context,
	token, password string,
) error {
	validPassword, err := domain.ValidatePassword(password)
	if err != nil {
		return err
	}

	validToken, err := domain.ValidateToken(token)
	if err != nil {
		return err
	}

	s, err := uc.repo.session.GetByToken(ctx, validToken)
	if err != nil {
		return err
	}

	if s.Type != domain.ResetPasswordToken || (s.ValidUntil != nil && s.ValidUntil.Before(time.Now())) {
		return domain.ErrInvalidToken
	}

	u, err := uc.repo.user.Get(ctx, s.UserId)
	if err != nil {
		return err
	}

	u.HashedPassword, err = domain.HashPassword(validPassword)
	if err != nil {
		return err
	}

	updates := domain.EntityUpdate{domain.UserHashedPasswordFieldName: u.HashedPassword}
	if u.EmailVerifyTime == nil {
		updates[domain.UserEmailVerifyTimeFieldName] = time.Now()
	}
	err = uc.repo.user.Update(ctx, u.Id, updates)
	if err != nil {
		return err
	}

	return nil
}

func (uc UserUseCase) SignOut(
	ctx context.Context,
	token string,
) error {
	s, err := uc.getGeneralValidSession(ctx, token)
	if err != nil {
		return err
	}

	err = uc.repo.session.Update(ctx, s.Id, domain.EntityUpdate{domain.SessionValidUntilFieldName: time.Now()})
	if err != nil {
		return err
	}

	return nil
}

func (uc UserUseCase) ChangePassword(
	ctx context.Context,
	token string,
	oldPassword, newPassword string,
) error {
	s, err := uc.getGeneralValidSession(ctx, token)
	if err != nil {
		return err
	}

	validOldPassword, err := domain.ValidatePassword(oldPassword)
	if err != nil {
		return err
	}

	validNewPassword, err := domain.ValidatePassword(newPassword)
	if err != nil {
		return err
	}

	u, err := uc.repo.user.Get(ctx, s.UserId)
	if err != nil {
		return err
	}

	if err = u.HashedPassword.Match(validOldPassword); err != nil {
		return domain.ErrInvalidPassword
	}

	u.HashedPassword, err = domain.HashPassword(validNewPassword)
	if err != nil {
		return err
	}
	err = uc.repo.user.Update(ctx, u.Id, domain.EntityUpdate{
		domain.UserHashedPasswordFieldName: u.HashedPassword,
	})
	if err != nil {
		return err
	}

	return nil
}

type ProfileDTO struct {
	Email           domain.Email
	EmailIsVerified bool
	Fullname        domain.Fullname
	Avatar          string
}

func (uc UserUseCase) GetProfile(
	ctx context.Context,
	token string,
) (p ProfileDTO, err error) {
	s, err := uc.getGeneralValidSession(ctx, token)
	if err != nil {
		return p, err
	}

	u, err := uc.repo.user.Get(ctx, s.UserId)
	if err != nil {
		return p, err
	}

	return ProfileDTO{
		Email:           u.Email,
		EmailIsVerified: u.EmailVerifyTime != nil,
		Fullname:        u.Fullname,
		Avatar:          u.Avatar,
	}, nil
}

type ProfileUpdateDTO struct {
	Fullname *string
	Avatar   *string
}

func (uc UserUseCase) UpdateProfile(
	ctx context.Context,
	token string,
	profileUpdate ProfileUpdateDTO,
) error {
	s, err := uc.getGeneralValidSession(ctx, token)
	if err != nil {
		return err
	}

	u, err := uc.repo.user.Get(ctx, s.UserId)
	if err != nil {
		return err
	}

	updates := domain.EntityUpdate{}
	{
		if profileUpdate.Fullname != nil {
			validFullname, err := domain.ValidateFullname(*profileUpdate.Fullname)
			if err != nil {
				return err
			}
			updates[domain.UserFullnameFieldName] = validFullname
		}

		if profileUpdate.Avatar != nil {
			// TODO check if avatar filename does exits then update it.
			updates[domain.UserAvatarFieldName] = *profileUpdate.Avatar
		}
	}
	err = uc.repo.user.Update(ctx, u.Id, updates)
	if err != nil {
		return err
	}

	return nil
}

// Upload
// Download
