package validator

import (
	"errors"
)

var (
	ErrPasswordTooShort = errors.New("password too short, password must be at least 8 characters")
)

func ValidatePassword(input string) error {
	if len(input) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}
