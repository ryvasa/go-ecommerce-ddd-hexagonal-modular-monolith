package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type Email struct {
	value string
}

func NewEmail(v string) (Email, error) {
	v = strings.TrimSpace(strings.ToLower(v))

	if !emailRegex.MatchString(v) {
		return Email{}, errors.New("invalid email format")
	}

	return Email{value: v}, nil
}

func (e Email) Value() string {
	return e.value
}
