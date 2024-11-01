package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/models"
)

const (
	minLength = 8

	alphabetLoverCase = "abcdefghijklmnopqrstuvwxyz"
	alphabetUpperCase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphabetNumber    = "0123456789"
	alphabetSymbol    = "!@#$%^&*()_+-=[]{};:'\",.<>/?|`~"
)

// Validator used to validate credentials.
// Must be created with New function.
type Validator struct {
	adminToken string
}

// New creates a new Validator instance.
// Requires a non-empty admin token.
// Returns a pointer to the Validator and an error if the token is empty.
func New(token string) (*Validator, error) {
	if token == "" {
		return nil, errors.New("empty admin token")
	}

	return &Validator{
		adminToken: token,
	}, nil
}

func (v *Validator) validateToken(token string) error {
	if token != v.adminToken {
		err := apperrors.ErrAccessDenied
		err.Message = "invalid admin token"
		return err
	}
	return nil
}

func (v *Validator) validateLogin(login string) error {
	if len(login) < minLength {
		err := apperrors.ErrWrongCredentials
		err.Message = fmt.Sprintf("login must be at least %d characters long", minLength)
		return err
	}

	alphabet := alphabetLoverCase + alphabetUpperCase + alphabetNumber

	for _, c := range login {
		if !strings.ContainsRune(alphabet, c) {
			err := apperrors.ErrWrongCredentials
			err.Message = fmt.Sprintf(
				"login must contain only letters and digits. Invalid character: %s",
				string(c),
			)
			return err
		}
	}
	return nil
}

func (v *Validator) validatePassword(pass string) error {
	if len(pass) < minLength {
		err := apperrors.ErrWrongCredentials
		err.Message = fmt.Sprintf("password must be at least %d characters long", minLength)
		return err
	}

	var (
		upperFound  bool
		lowerFound  bool
		numberFound bool
		symbolFound bool
	)

	for _, c := range pass {
		switch {
		case strings.ContainsRune(alphabetLoverCase, c):
			lowerFound = true
		case strings.ContainsRune(alphabetUpperCase, c):
			upperFound = true
		case strings.ContainsRune(alphabetNumber, c):
			numberFound = true
		case strings.ContainsRune(alphabetSymbol, c):
			symbolFound = true
		default:
			err := apperrors.ErrWrongCredentials
			err.Message = fmt.Sprintf(
				"password must contain only letters, digits and non space symbols. Unsupportded character: %s",
				string(c),
			)
			return err
		}
	}

	var errorMessages []string
	if !upperFound {
		errorMessages = append(errorMessages, "Must contains upper case letter.")
	}

	if !lowerFound {
		errorMessages = append(errorMessages, "Must contains lower case letter.")
	}

	if !numberFound {
		errorMessages = append(errorMessages, "Must contains digit.")
	}

	if !symbolFound {
		errorMessages = append(errorMessages, "Must contains non space symbol.")
	}

	if len(errorMessages) > 0 {
		err := apperrors.ErrWrongCredentials
		err.Message = fmt.Sprintf(
			"password validation failed: %s",
			strings.Join(errorMessages, " "),
		)
		return err
	}

	return nil
}

// ValidateCredentials validates models.Credentials.
// It checks if token is valid, login
// and password is valid according to the rules.
// If any of the checks fail, it returns an error.
// Otherwise, it returns nil.
func (v *Validator) ValidateCredentials(credentials models.Credentials) error {
	if err := v.validateToken(credentials.Token); err != nil {
		return err
	}
	if err := v.validateLogin(credentials.Login); err != nil {
		return err
	}
	if err := v.validatePassword(credentials.Password); err != nil {
		return err
	}
	return nil
}
