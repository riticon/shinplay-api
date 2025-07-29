// Package publicid provides public ID values in the same format as
// PlanetScale’s Rails application.
package publicid

import (
	"strings"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pkg/errors"
)

// Fixed nanoid parameters used in the Rails application.
const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	length   = 12
)

func Numberic() string {
	return "0123456789"
}

// New generates a unique public ID.
func New() (string, error) {
	return nanoid.Generate(alphabet, length)
}

// Must is the same as New, but panics on error.
func Must() string {
	return nanoid.MustGenerate(alphabet, length)
}

func MustWith(length int, alphabet string) string {
	if length <= 0 {
		panic("length must be greater than 0")
	}
	return nanoid.MustGenerate(alphabet, length)
}

// Validate checks if a given field name’s public ID value is valid according to
// the constraints defined by package publicid.
func Validate(fieldName, id string) error {
	if id == "" {
		return errors.Errorf("%s cannot be blank", fieldName)
	}

	if len(id) != length {
		return errors.Errorf("%s should be %d characters long", fieldName, length)
	}

	if strings.Trim(id, alphabet) != "" {
		return errors.Errorf("%s has invalid characters", fieldName)
	}

	return nil
}
