package hw09structvalidator

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

func TestSliceValidator(t *testing.T) {
	type v struct {
		IntSlice   []int     `validate:"min:10|max:20"`
		StrSlice   []string  `validate:"len:5"`
		FloatSlice []float64 `validate:"min:10.0"`
	}

	cases := []TestCase{
		{
			Name:  "int slice",
			Param: v{IntSlice: []int{1, 12, 21}},
			ExpectedErr: ValidationErrors{
				ValidationError{Field: "IntSlice", Err: ErrValidationMin},
				ValidationError{Field: "IntSlice", Err: ErrValidationMax},
			},
		},
		{
			Name:  "string slice",
			Param: v{StrSlice: []string{"123", "123456"}},
			ExpectedErr: ValidationErrors{
				ValidationError{Field: "StrSlice", Err: ErrValidationLen},
			},
		},
		{
			Name:  "unsupported type slice",
			Param: v{FloatSlice: []float64{1.0, 2.0}},
			ExpectedErr: ValidationErrors{
				ValidationError{Field: "FloatSlice", Err: ErrValidationType},
			},
		},
	}

	for _, c := range cases {
		t.Log(c.Name)
		err := Validate(c.Param)
		require.Equal(t, err, c.ExpectedErr)
	}
}
