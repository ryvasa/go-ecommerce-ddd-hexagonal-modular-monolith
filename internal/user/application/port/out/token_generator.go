package out

type TokenGenerator interface {
	Generate(userID string, roles []string) (string, error)
}
