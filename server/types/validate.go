package types

import (
	"fmt"
	"github.com/mccoyst/validate"
	"regexp"
	"unicode/utf8"
)

const (
	MIN_PASS_LENGTH   = 8
	MAX_HANDLE_LENGTH = 16
)

// Regexes
var (
	handleRegex = regexp.MustCompile(`^\p{L}[\d\p{L}]*`)
)

func NewValidator() *validate.V {
	vd := make(validate.V)
	vd["handle"] = validateHandle
	return &vd
}

func validateHandle(i interface{}) error {
	handle := i.(string)
	if handle == "" {
		return fmt.Errorf("Handle is a required field for signup")
	} else if utf8.RuneCountInString(handle) > MAX_HANDLE_LENGTH {
		return fmt.Errorf("Handle is too long, max length is %d", MAX_HANDLE_LENGTH)
	} else if !handleRegex.MatchString(handle) {
		return fmt.Errorf(handle + " contains illegal characters")
	} else {
		return nil
	}
}
