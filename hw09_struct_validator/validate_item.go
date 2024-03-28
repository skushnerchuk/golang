package hw09structvalidator

import (
	"fmt"
)

type Rule struct {
	Operation string
	Value     string
}

func (v *Rule) String() string {
	return fmt.Sprintf("%s:%s", v.Operation, v.Value)
}
