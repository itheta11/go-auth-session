package utils

import (
	"errors"
	"regexp"
	"strings"
)

func ValidatePassword(password, username, email string) error {
	// Length check
	if len(password) < 12 {
		return errors.New("password must be at least 12 characters long")
	}
	if len(password) > 64 {
		return errors.New("password must not exceed 64 characters")
	}

	// Uppercase
	if match, _ := regexp.MatchString(`[A-Z]`, password); !match {
		return errors.New("password must contain at least one uppercase letter")
	}

	// Lowercase
	if match, _ := regexp.MatchString(`[a-z]`, password); !match {
		return errors.New("password must contain at least one lowercase letter")
	}

	// Digit
	if match, _ := regexp.MatchString(`[0-9]`, password); !match {
		return errors.New("password must contain at least one digit")
	}

	// Special char
	if match, _ := regexp.MatchString(`[\W_]`, password); !match {
		return errors.New("password must contain at least one special character")
	}

	// No spaces
	if strings.Contains(password, " ") {
		return errors.New("password must not contain spaces")
	}

	// Prevent username or email reuse
	if username != "" && strings.Contains(strings.ToLower(password), strings.ToLower(username)) {
		return errors.New("password must not contain username")
	}
	if email != "" && strings.Contains(strings.ToLower(password), strings.ToLower(strings.Split(email, "@")[0])) {
		return errors.New("password must not contain parts of email")
	}

	// Basic blacklist (optional)
	blacklist := []string{"password", "123456", "qwerty", "letmein"}
	for _, bad := range blacklist {
		if strings.EqualFold(password, bad) {
			return errors.New("password is too common")
		}
	}

	return nil
}
