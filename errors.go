package fyneloader

import "fmt"

// ConstError is a simple constant error type.
type ConstError string

func (e ConstError) Error() string {
	return string(e)
}

const (
	// ErrNoWidgetType indicates that an element map did not have a type field.
	ErrNoWidgetType ConstError = "no type tag on element definition"

	// ErrUnknownFileExt indicates that a file path did not have a known
	// extension.
	ErrUnknownFileExt ConstError = "unknown file extension"
)

// FunctionTypeError is an error which indicates that a function type did not
// match any of the allowed types.
type FunctionTypeError struct {
	Func interface{}
}

func (e FunctionTypeError) Error() string {
	return fmt.Sprintf("invalid function type %T", e.Func)
}

// UndefinedFunctionError is an error which indicates that the function with
// the given name was not registered.
type UndefinedFunctionError struct {
	Name string
}

func (e UndefinedFunctionError) Error() string {
	return fmt.Sprintf("no function %q defined", e.Name)
}

// UnknownElementType is an error which indicates that the element with the
// given name is unknown.
type UnknownElementType struct {
	TypeName string
}

func (e UnknownElementType) Error() string {
	return fmt.Sprintf("unknown element %q", e.TypeName)
}
