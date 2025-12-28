package security

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (b *BcryptHasher) Hash(p string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(bytes), err
}

func (b *BcryptHasher) Compare(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(plain),
	) == nil
}
