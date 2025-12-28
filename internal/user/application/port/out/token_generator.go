package out

type TokenGenerator interface {
	Generate(userID string, role string) (string, error)
}
