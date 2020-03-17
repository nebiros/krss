package db

import (
	"fmt"
	"strings"
)

type AssignmentsBuilder struct {
	assignments []string
	values      []interface{}
}

func (s *AssignmentsBuilder) Add(key string, value interface{}) {
	if s.assignments == nil {
		s.assignments = make([]string, 0, 1)
	}
	if s.values == nil {
		s.values = make([]interface{}, 0, 1)
	}
	s.assignments = append(s.assignments, fmt.Sprintf("%s = ?", key))
	s.values = append(s.values, value)
}

func (s AssignmentsBuilder) Assignments() string {
	return strings.Join(s.assignments, ", ")
}

func (s AssignmentsBuilder) Values() []interface{} {
	return s.values
}
