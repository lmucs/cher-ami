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
	handleRegex        = regexp.MustCompile(`^[\p{L}\p{M}][\d\p{L}\p{M}]*$`)
	emailRegex         = regexp.MustCompile(`^\w+([-+.']\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)
	userAttributeRegex = regexp.MustCompile(`^((fir|la)stname|gender|birthday|bio|interests|languages|location)$`)
)

func NewValidator() *validate.V {
	vd := validate.V{
		"handle":          validateHandle,
		"email":           validateEmail,
		"password":        validatePassword,
		"messageop":       validateMessageOp,
		"messageresource": validateMessageResource,
		"messagevalue":    validateMessageValue,
		"userresource":    validateUserResource,
		"uservalue":       validateUserValue,
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
		return fmt.Errorf(email + " is an invalid email")
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

//
// EditMessage Validation
//

func validateMessageOp(i interface{}) error {
	op := i.(string)
	if op == "" {
		return fmt.Errorf("Required field for message patch")
	} else if op != "update" && op != "publish" && op != "unpublish" {
		return fmt.Errorf(op)
	} else {
		return nil
	}
}

func validateMessageResource(i interface{}) error {
	resource := i.(string)
	if resource == "" {
		return fmt.Errorf("Required field for message patch")
	} else if resource != "content" && resource != "image" && resource != "circle" {
		return fmt.Errorf(resource)
	} else {
		return nil
	}
}

func validateMessageValue(i interface{}) error {
	value := i.(string)
	if value == "" {
		return fmt.Errorf("Required field for message patch")
	} else {
		return nil
	}
}

func validateUserResource(i interface{}) error {
	resource := i.(string)
	if resource == "" {
		return fmt.Errorf("Required field for user patch")
		// !r.MatchString(resource)
	} else if resource != "firstname" && resource != "lastname" && resource != "gender" && resource != "birthday" && resource != "bio" && resource != "interests" && resource != "languages" && resource != "location" {
		return fmt.Errorf(resource)
	} else {
		return nil
	}
}

func validateUserValue(i interface{}) error {
	value := i.(string)
	if value == "" {
		return fmt.Errorf("Required field for user patch")
	} else {
		return nil
	}
}
