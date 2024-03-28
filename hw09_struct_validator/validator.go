package hw09structvalidator

import (
	"reflect"
	"strings"
	"unicode"
)

const (
	Len    = "len"
	Min    = "min"
	Max    = "max"
	In     = "in"
	Regexp = "regexp"
)

func createRule(rule string) (*Rule, error) {
	condition := strings.Split(rule, ":")
	if len(condition) != 2 || condition[0] == "" || condition[1] == "" {
		return nil, ErrValidationUnsupportedRule
	}
	return &Rule{
		Operation: condition[0],
		Value:     condition[1],
	}, nil
}

func parseTag(tag string) ([]Rule, error) {
	rules := make([]Rule, 0)
	for _, r := range strings.Split(tag, "|") {
		rule, err := createRule(r)
		if err != nil {
			return nil, ErrValidationUnsupportedRule
		}
		rules = append(rules, *rule)
	}
	return rules, nil
}

func validateNeeded(field reflect.StructField) bool {
	name := field.Name
	tag := field.Tag.Get("validate")
	return tag != "" && !field.Anonymous && unicode.IsUpper(rune(name[0]))
}

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)

	if value.Kind() == reflect.Ptr && !value.IsNil() {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	var errList ValidationErrors

	for _, field := range reflect.VisibleFields(value.Type()) {
		if !validateNeeded(field) {
			continue
		}
		name := field.Name
		tag := field.Tag.Get("validate")
		if err := validateField(value.FieldByName(name), name, tag); err != nil {
			errList = append(errList, err...)
		}
	}

	if len(errList) == 0 {
		return nil
	}

	return errList
}

func validateField(field reflect.Value, name string, tag string) []ValidationError {
	errors := make([]ValidationError, 0)

	rules, err := parseTag(tag)
	if err != nil {
		errors = append(errors, ValidationError{Field: name, Err: err})
		return errors
	}

	switch field.Kind() { //nolint:exhaustive
	case reflect.Int:
		return validateInt(field.Int(), name, rules)
	case reflect.String:
		return validateString(field.String(), name, rules)
	case reflect.Slice:
		slice := make([]interface{}, field.Len())
		for i := 0; i < field.Len(); i++ {
			slice[i] = field.Index(i).Interface()
		}
		return validateSlice(slice, name, rules)
	default:
		errors = append(errors, ValidationError{Field: name, Err: ErrValidationType})
	}
	return errors
}
