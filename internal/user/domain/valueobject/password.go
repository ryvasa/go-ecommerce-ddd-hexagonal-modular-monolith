package valueobject

import "errors"

type Password struct {
	hash string
}

func NewPasswordFromPlain(
	plain string,
	hasher PasswordHasher,
) (Password, error) {
	if err := validatePasswordPolicy(plain); err != nil {
		return Password{}, err
	}

	hash, err := hasher.Hash(plain)
	if err != nil {
		return Password{}, err
	}

	return Password{hash: hash}, nil
}

func NewHashedPassword(hash string) (Password, error) {
	if hash == "" {
		return Password{}, errors.New("password hash cannot be empty")
	}
	return Password{hash: hash}, nil
}

func (p Password) Hash() string {
	return p.hash
}

func (p Password) Verify(plain string, hasher PasswordHasher) bool {
	return hasher.Compare(p.hash, plain)
}
