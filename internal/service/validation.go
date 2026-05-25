package service

import (
	"net/mail"
	"strings"
)

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return err
	}
	return nil
}
