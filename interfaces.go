package logr

import (
	"fmt"
	"strings"
)

type Interfaces []any

func (i Interfaces) Strings() []string {
	var r []string
	for _, v := range i {
		r = append(r, fmt.Sprintf("%v", v))
	}
	return r
}

// SSV creates a space separated values string
func (i Interfaces) SSV() string {
	return strings.TrimSpace(strings.Join(i.Strings(), " "))
}
