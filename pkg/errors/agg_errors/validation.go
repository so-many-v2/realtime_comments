package agg_errors

import "fmt"

type ValidationError struct {
	Field string
	Msg   string
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("%s : %s", ve.Field, ve.Msg)
}
