package out

// TokenGenerator creates JWT tokens
type TokenGenerator interface {
	Generate(userID string, roles []string) (string, error)
}
