package validation

import (
	"strings"
	"unicode"
)

var (
	validMails = map[string]struct{}{
		"mail":       {},
		"gmail":      {},
		"yahoo":      {},
		"outlook":    {},
		"yandex":     {},
		"zoho":       {},
		"protonmail": {},
		"icloud":     {},
		"aol":        {},
		"gmx":        {},
	}

	validDomens = map[string]struct{}{
		"com": {},
		"ru":  {},
	}
)

const (
	minLenNickname = 5
	maxLenNickname = 15

	minLenPassword = 6
	maxLenPassword = 64

	specialChars = `!@#$%^&*()_+={}[\\]:;\"'<>,.?~\\-`
)

func IsValidNickname(nickname string) bool {
	if len(nickname) < minLenNickname ||
		len(nickname) > maxLenNickname {
		return false
	}

	for _, char := range nickname {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}

	return true
}

func IsValidEmail(email string) bool {
	if len(email) > 46 {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	if len(parts[0]) < 5 || len(parts[0]) > 31 {
		return false
	}

	for _, char := range parts[0] {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}

	domainParts := strings.Split(parts[1], ".")

	if len(domainParts) != 2 {
		return false
	}

	if _, ok := validMails[domainParts[0]]; !ok {
		return false
	}

	if _, ok := validDomens[domainParts[1]]; !ok {
		return false
	}

	return true
}

func IsValidNewPassword(password, PasswordConfirm string) bool {
	if password != PasswordConfirm {
		return false
	}

	if len(password) < minLenPassword ||
		len(password) > maxLenPassword {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else if contains(specialChars, char) {
			hasSpecial = true
		} else {
			return false
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}

func IsValidPassword(password string) bool {
	if len(password) < minLenPassword ||
		len(password) > maxLenPassword {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else if contains(specialChars, char) {
			hasSpecial = true
		} else {
			return false
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}

func contains(specialChars string, char rune) bool {
	for _, specialChar := range specialChars {
		if char == specialChar {
			return true
		}
	}
	return false
}
