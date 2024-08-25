package utils

import (
	"net/mail"
	"strings"
)

const (
	ErrTooShort         = "Password is too short"
	ErrTooLong          = "Password is too long"
	ErrWrongFormat      = "Password does not contain required symbols"
	ErrWrongEmailFormat = "Wrong email format"
	ErrDoNotMatch       = "Passwords do not match"
)

const (
	minLen = 8
	maxLen = 32
)

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}

func validatePassword(password string) []string {
	var errors []string

	switch {
	case len(password) < minLen:
		errors = append(errors, ErrTooShort)
	case len(password) > maxLen:
		errors = append(errors, ErrTooLong)
	case !checkPassword(password):
		errors = append(errors, ErrWrongFormat)
	}

	return errors
}

func checkPassword(password string) bool {
	specialChars := "!@#$%^&*?"

	for _, char := range password {
		if strings.ContainsRune(specialChars, char) {
			return true
		}
	}

	return false
}

func Validate(email, password, passwordRepeat string) []string {
	errors := validatePassword(password)

	if password != passwordRepeat {
		errors = append(errors, ErrDoNotMatch)
	}

	if !validateEmail(email) {
		errors = append(errors, ErrWrongEmailFormat)
	}

	return errors
}
