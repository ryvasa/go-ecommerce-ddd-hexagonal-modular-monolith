package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
)

type User struct {
	id              string
	email           valueobject.Email
	password        valueobject.Password
	roles           map[string]struct{}
	firstName       string
	lastName        string
	username        string
	phone           string
	birthDate       time.Time
	emailVerifiedAt time.Time
}

func RehydrateUser(
	id string,
	email valueobject.Email,
	password valueobject.Password,
) *User {
	return &User{
		id:       id,
		email:    email,
		password: password,
		// active:   active,
		// roles:    roleMap,
	}
}

func (u *User) ChangePassword(
	oldPlain string,
	newPlain string,
	hasher valueobject.PasswordHasher,
) error {

	if !u.password.Verify(oldPlain, hasher) {
		return errors.New("invalid current password")
	}

	newPassword, err := valueobject.NewPasswordFromPlain(newPlain, hasher)
	if err != nil {
		return err
	}

	u.password = newPassword
	return nil
}

func (u *User) AssignRole(role string) error {
	if _, exists := u.roles[role]; exists {
		return errors.New("role already assigned")
	}
	u.roles[role] = struct{}{}
	return nil
}

func NewUser(
	email valueobject.Email,
	password valueobject.Password,
	username string,
	firstName string,
	lastName string,
) *User {

	return &User{
		id:        uuid.NewString(),
		email:     email,
		username:  username,
		firstName: firstName,
		lastName:  lastName,
		password:  password,
		roles:     map[string]struct{}{"user": {}},
	}
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Email() valueobject.Email {
	return u.email
}

func (u *User) Password() valueobject.Password {
	return u.password
}

func (u *User) Roles() []string {
	roles := make([]string, 0, len(u.roles))
	for r := range u.roles {
		roles = append(roles, r)
	}
	return roles
}

func (u *User) Username() string {
	return u.username
}

func (u *User) FirstName() string {
	return u.firstName
}

func (u *User) LastName() string {
	return u.lastName
}
