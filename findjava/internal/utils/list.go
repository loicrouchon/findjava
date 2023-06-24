package utils

import "strings"

type List []string

func (i *List) String() string {
	return "[" + strings.Join(*i, ", ") + "]"
}

func (i *List) Set(value string) error {
	*i = append(*i, value)
	return nil
}
