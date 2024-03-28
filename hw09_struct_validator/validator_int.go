package hw09structvalidator

import (
	"slices"
	"strconv"
	"strings"
)

func validateInt(fieldValue int64, name string, rules []Rule) []ValidationError {
	errors := make([]ValidationError, 0)

	for _, r := range rules {
		switch r.Operation {
		case Min:
			v, err := strconv.ParseInt(r.Value, 10, 64)
			if err != nil {
				errors = append(errors, ValidationError{Field: name, Err: err})
				continue
			}
			if fieldValue < v {
				errors = append(errors, ValidationError{Field: name, Err: ErrValidationMin})
			}
		case Max:
			v, err := strconv.ParseInt(r.Value, 10, 64)
			if err != nil {
				errors = append(errors, ValidationError{Field: name, Err: err})
				continue
			}
			if fieldValue > v {
				errors = append(errors, ValidationError{Field: name, Err: ErrValidationMax})
			}
		case In:
			list := strings.Split(r.Value, ",")
			if !slices.Contains(list, strconv.FormatInt(fieldValue, 10)) {
				errors = append(errors, ValidationError{Field: name, Err: ErrValidationIn})
			}
		default:
			errors = append(errors, ValidationError{Field: name, Err: ErrValidationUnsupportedRule})
		}
	}

	return errors
}
