package hw09structvalidator

func validateSlice(slice []interface{}, name string, rules []Rule) []ValidationError {
	errors := make([]ValidationError, 0)
	for _, value := range slice {
		switch v := value.(type) {
		case string:
			errors = append(errors, validateString(v, name, rules)...)
		case int:
			errors = append(errors, validateInt(int64(v), name, rules)...)
		default:

			errors = append(errors, ValidationError{Field: name, Err: ErrValidationType})
			return errors
		}
	}
	return errors
}
