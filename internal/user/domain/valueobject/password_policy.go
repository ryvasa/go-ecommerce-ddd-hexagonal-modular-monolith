package valueobject

import "errors"

func validatePasswordPolicy(p string) error {
	if len(p) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}
