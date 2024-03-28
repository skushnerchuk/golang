package hw09structvalidator

import (
	"regexp/syntax"
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

func TestStringLenValidator(t *testing.T) {
	type v struct {
		Version string `validate:"len:5"`
	}

	cases := []TestCase{
		{
			Name:        "too long",
			Param:       v{Version: "0.0.11"},
			ExpectedErr: ValidationErrors{ValidationError{Field: "Version", Err: ErrValidationLen}},
		},
		{
			Name:        "no error",
			Param:       &v{Version: "0.1"},
			ExpectedErr: nil,
		},
	}

	for _, c := range cases {
		t.Log(c.Name)
		err := Validate(c.Param)
		require.Equal(t, err, c.ExpectedErr)
	}
}

func TestStringRegexpValidator(t *testing.T) {
	type correctRegexp struct {
		Email string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	}

	type incorrectRegexp struct {
		Email string `validate:"regexp:^$["`
	}

	cases := []TestCase{
		{
			Name:        "correct regexp but incorrect value",
			Param:       correctRegexp{Email: "user@"},
			ExpectedErr: ValidationErrors{ValidationError{Field: "Email", Err: ErrValidationRegexp}},
		},
		{
			Name:  "incorrect regexp but correct value",
			Param: incorrectRegexp{Email: "user@example.com"},
			ExpectedErr: ValidationErrors{
				ValidationError{
					Field: "Email",
					Err:   &syntax.Error{Code: "missing closing ]", Expr: "["},
				},
			},
		},
		{
			Name:        "no error",
			Param:       &correctRegexp{Email: "user@example.com"},
			ExpectedErr: nil,
		},
	}

	for _, c := range cases {
		t.Log(c.Name)
		err := Validate(c.Param)
		require.Equal(t, err, c.ExpectedErr)
	}
}

func TestStringInValidator(t *testing.T) {
	type v struct {
		Role string `validate:"in:admin,stuff"`
	}

	cases := []TestCase{
		{
			Name:        "in error",
			Param:       v{Role: "user"},
			ExpectedErr: ValidationErrors{ValidationError{Field: "Role", Err: ErrValidationIn}},
		},
		{
			Name:        "no error",
			Param:       v{Role: "admin"},
			ExpectedErr: nil,
		},
	}

	for _, c := range cases {
		t.Log(c.Name)
		err := Validate(c.Param)
		require.Equal(t, err, c.ExpectedErr)
	}
}
