package out

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) bool
}
