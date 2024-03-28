package hw09structvalidator

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

func TestIntMinMaxValidator(t *testing.T) {
	type v struct {
		Age int `validate:"min:18|max:50"`
	}

	cases := []TestCase{
		{
			Name: "no error", Param: v{Age: 22}, ExpectedErr: nil,
		},
		{
			Name:        "min rule",
			Param:       v{Age: 16},
			ExpectedErr: ValidationErrors{ValidationError{Field: "Age", Err: ErrValidationMin}},
		},
		{
			Name:        "max rule",
			Param:       v{Age: 60},
			ExpectedErr: ValidationErrors{ValidationError{Field: "Age", Err: ErrValidationMax}},
		},
	}

	for _, c := range cases {
		t.Log(c.Name)
		err := Validate(c.Param)
		require.Equal(t, err, c.ExpectedErr)
	}
}

func TestIntInValidator(t *testing.T) {
	type v struct {
		Code int `validate:"in:200,404,500"`
	}

	cases := []TestCase{
		{Name: "no error", Param: v{Code: 200}, ExpectedErr: nil},
		{
			Name:        "not in",
			Param:       v{Code: 301},
			ExpectedErr: ValidationErrors{ValidationError{Field: "Code", Err: ErrValidationIn}},
		},
	}

	for _, c := range cases {
		t.Logf(c.Name)
		err := Validate(c.Param)
		require.Equal(t, err, c.ExpectedErr)
	}
}

func TestIntIncorrectValueValidator(t *testing.T) {
	type v struct {
		Code int `validate:"min:error"`
	}

	cases := []TestCase{
		{
			Name:  "string instead number",
			Param: v{Code: 200},
			ExpectedErr: ValidationErrors{
				ValidationError{
					Field: "Code", Err: &strconv.NumError{
						Func: "ParseInt",
						Num:  "error",
						Err:  strconv.ErrSyntax,
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Logf(c.Name)
		err := Validate(c.Param)
		require.Equal(t, err, c.ExpectedErr)
	}
}
