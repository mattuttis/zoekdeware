package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidEmail = errors.New("invalid email format")
	emailRegex      = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type Email string

func NewEmail(value string) (Email, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if !emailRegex.MatchString(normalized) {
		return "", ErrInvalidEmail
	}
	return Email(normalized), nil
}

func (e Email) String() string {
	return string(e)
}

func (e Email) Domain() string {
	parts := strings.Split(string(e), "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}
