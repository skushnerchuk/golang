package hw09structvalidator

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotStruct                 = errors.New("not a struct")
	ErrValidationType            = errors.New("unsupported type")
	ErrValidationMin             = errors.New("rule min")
	ErrValidationMax             = errors.New("rule max")
	ErrValidationLen             = errors.New("rule len")
	ErrValidationIn              = errors.New("rule in")
	ErrValidationRegexp          = errors.New("rule regexp")
	ErrValidationUnsupportedRule = errors.New("unsupported rule")
)

type ValidationError struct {
	Field string
	Err   error
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("Field %s: %s", v.Field, v.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	errText := strings.Builder{}
	for _, err := range v {
		errText.WriteString(err.Error())
	}
	return errText.String()
}
