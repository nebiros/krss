package model

import "fmt"

type ErrEmptyArgument struct {
	Name  string
	Value interface{}
}

func (e ErrEmptyArgument) Error() string {
	return fmt.Sprintf("argument '%s' seems empty: '%v'", e.Name, e.Value)
}
