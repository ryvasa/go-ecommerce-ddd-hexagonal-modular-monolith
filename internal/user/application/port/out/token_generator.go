package out

type TokenGenerator interface {
	Generate(userID string) (string, error)
}
