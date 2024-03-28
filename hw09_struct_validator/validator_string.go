package hw09structvalidator

import (
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func validateString(fieldValue, name string, rules []Rule) []ValidationError {
	errors := make([]ValidationError, 0)

	for _, r := range rules {
		switch r.Operation {
		case Len:
			length, err := strconv.Atoi(r.Value)
			if err != nil {
				errors = append(errors, ValidationError{Field: name, Err: err})
				continue
			}
			if len(fieldValue) > length {
				errors = append(errors, ValidationError{Field: name, Err: ErrValidationLen})
			}
		case In:
			list := strings.Split(r.Value, ",")
			if !slices.Contains(list, fieldValue) {
				errors = append(errors, ValidationError{Field: name, Err: ErrValidationIn})
			}
		case Regexp:
			re, err := regexp.Compile(r.Value)
			if err != nil {
				errors = append(errors, ValidationError{Field: name, Err: err})
				continue
			}
			match := re.MatchString(fieldValue)
			if !match {
				errors = append(errors, ValidationError{Field: name, Err: ErrValidationRegexp})
			}
		default:
			errors = append(errors, ValidationError{Field: name, Err: ErrValidationUnsupportedRule})
		}
	}
	return errors
}
