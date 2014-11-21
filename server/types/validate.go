package types

import (
	"fmt"
	"github.com/mccoyst/validate"
	"regexp"
	"unicode/utf8"
)

const (
	MIN_PASS_LENGTH   = 8
	MAX_PASS_LENGTH   = 50
	MAX_HANDLE_LENGTH = 16
)

// Regexes
var (
	handleRegex = regexp.MustCompile(`^[\p{L}\p{M}][\d\p{L}\p{M}]*$`)
	emailRegex  = regexp.MustCompile(`^\w+([-+.']\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)
)

func NewValidator() *validate.V {
	vd := validate.V{
		"handle":   validateHandle,
		"email":    validateEmail,
		"password": validatePassword,
	}
	return &vd
}

func validateHandle(i interface{}) error {
	handle := i.(string)
	if handle == "" {
		return fmt.Errorf("Required field for signup")
	} else if utf8.RuneCountInString(handle) > MAX_HANDLE_LENGTH {
		return fmt.Errorf("Too long, max length is %d", MAX_HANDLE_LENGTH)
	} else if !handleRegex.MatchString(handle) {
		return fmt.Errorf(handle + " contains illegal characters")
	} else {
		return nil
	}
}

func validateEmail(i interface{}) error {
	email := i.(string)
	if email == "" {
		return fmt.Errorf("Required field for signup")
	} else if !emailRegex.MatchString(email) {
		return fmt.Errorf(email + " is an invalid email.")
	} else {
		return nil
	}
}

func validatePassword(i interface{}) error {
	password := i.(string)
	passwordLen := utf8.RuneCountInString(password)
	if password == "" {
		return fmt.Errorf("Required field for signup")
	} else if passwordLen < MIN_PASS_LENGTH {
		return fmt.Errorf("Too short, minimum length is %d", MIN_PASS_LENGTH)
	} else if passwordLen > MAX_PASS_LENGTH {
		return fmt.Errorf("Too long, maximum length is %d", MAX_PASS_LENGTH)
	} else {
		return nil
	}
}
