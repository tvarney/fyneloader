package fyneloader

import (
	"fmt"
	"strings"
)

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

	// ErrInvalidOption indicates that an option given was not valid.
	ErrInvalidOption ConstError = "invalid option"
)

// ArrayIndexOutOfBoundsError is an error indicating that the given index was
// out of bounds for the source array.
type ArrayIndexOutOfBoundsError struct {
	Index int64
}

func (e ArrayIndexOutOfBoundsError) Error() string {
	return fmt.Sprintf("array index %d out of bounds", e.Index)
}

// ConflictingKeysError is an error which indicates that some keys in the
// configuration conflict.
type ConflictingKeysError struct {
	Keys []string
}

func (e ConflictingKeysError) Error() string {
	if len(e.Keys) == 0 {
		return "conflicting keys"
	}
	builder := &strings.Builder{}
	fmt.Fprintf(builder, "conflicting keys: %s", e.Keys[0])
	for _, v := range e.Keys[1:] {
		fmt.Fprintf(builder, ", %s", v)
	}
	return builder.String()
}

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
