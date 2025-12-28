package entity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
)

type User struct {
	id       string
	email    valueobject.Email
	password valueobject.Password
	active   bool
	roles    map[string]struct{}
}

func RehydrateUser(
	id string,
	email valueobject.Email,
	password valueobject.Password,
) *User {

	// roleMap := make(map[string]struct{}, len(roles))
	// for _, r := range roles {
	// 	roleMap[r] = struct{}{}
	// }

	return &User{
		id:       id,
		email:    email,
		password: password,
		// active:   active,
		// roles:    roleMap,
	}
}

func (u *User) Activate() error {
	if u.active {
		return errors.New("user already active")
	}
	u.active = true
	return nil
}

func (u *User) Deactivate() error {
	if !u.active {
		return errors.New("user already inactive")
	}
	u.active = false
	return nil
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
) *User {

	return &User{
		id:       uuid.NewString(),
		email:    email,
		password: password,
		active:   true,
		roles:    map[string]struct{}{"user": {}},
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

func (u *User) IsActive() bool {
	return u.active
}

func (u *User) Roles() []string {
	roles := make([]string, 0, len(u.roles))
	for r := range u.roles {
		roles = append(roles, r)
	}
	return roles
}
