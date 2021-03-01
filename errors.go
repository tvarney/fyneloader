package fyneloader

import "fmt"

// FunctionTypeError is an error which indicates that a function type did not
// match any of the allowed types.
type FunctionTypeError struct {
	Func interface{}
}

func (e FunctionTypeError) Error() string {
	return fmt.Sprintf("invalid function type %T", e.Func)
}
