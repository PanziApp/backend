package domain

import "time"

type User struct {
	Id              EntityId
	CreateTime      time.Time
	Email           Email
	EmailVerifyTime *time.Time
	HashedPassword  HashedPassword
	Fullname        Fullname
	Avatar          string
}

const (
	UserEmailVerifyTimeFieldName EntityFieldName = "user_email_verify_time"
	UserHashedPasswordFieldName  EntityFieldName = "user_hashed_password"
	UserFullnameFieldName        EntityFieldName = "user_fullname"
	UserAvatarFieldName          EntityFieldName = "user_avatar"
)
