package hw09structvalidator

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

type TestCase struct {
	Name        string
	Param       interface{}
	ExpectedErr error
}

func TestCheckParamIsStructure(t *testing.T) {
	type v struct {
		SameField string
	}

	cases := []TestCase{
		{Name: "struct by value", Param: v{SameField: ""}, ExpectedErr: nil},
		{Name: "struct by pointer", Param: &v{SameField: ""}, ExpectedErr: nil},
		{Name: "number instead struct", Param: 5, ExpectedErr: ErrNotStruct},
		{Name: "<nil> instead struct", Param: nil, ExpectedErr: ErrNotStruct},
	}

	for _, c := range cases {
		t.Log(c.Name)
		err := Validate(c.Param)
		require.Equal(t, err, c.ExpectedErr)
	}
}

func TestUnsupportedRules(t *testing.T) {
	type v struct {
		SameField string `validate:"required:10"`
	}

	cases := []TestCase{
		{
			Name:  "unsupported rule",
			Param: v{SameField: ""},
			ExpectedErr: ValidationErrors{
				ValidationError{Field: "SameField", Err: ErrValidationUnsupportedRule},
			},
		},
	}

	for _, c := range cases {
		t.Log(c.Name)
		err := Validate(c.Param)
		require.Equal(t, err, c.ExpectedErr)
	}
}

func TestUnsupportedType(t *testing.T) {
	type v struct {
		SameField float64 `validate:"min:10.55"`
	}

	cases := []TestCase{
		{
			Name:  "unsupported type",
			Param: v{SameField: 9.0},
			ExpectedErr: ValidationErrors{
				ValidationError{Field: "SameField", Err: ErrValidationType},
			},
		},
	}

	for _, c := range cases {
		t.Log(c.Name)
		err := Validate(c.Param)
		require.Equal(t, err, c.ExpectedErr)
	}
}
