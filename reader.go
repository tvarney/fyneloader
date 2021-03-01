package fyneloader

import (
	"reflect"

	"fyne.io/fyne/v2"
)

// CreateElementFn is the function type used for the element creation
// callbacks.
type CreateElementFn func(interface{}) fyne.CanvasObject

// Loader allows for loading UI definitions at runtime.
type Loader struct {
	callbacks map[string]interface{}
	elements  map[string]CreateElementFn
}

// New returns a new Loader instance.
func New() *Loader {
	return &Loader{
		callbacks: map[string]interface{}{},
		elements:  map[string]CreateElementFn{},
	}
}

// RegisterElement registers a new element callback.
//
// If the function callback is nil and there already exists a callback for the
// given name, the callback will be removed. This function will replace a
// callback with no error if a name is repeated.
func (l *Loader) RegisterElement(name string, fn CreateElementFn) {
	if fn == nil {
		_, ok := l.elements[name]
		if ok {
			delete(l.elements, name)
		}
		return
	}
	l.elements[name] = fn
}

// RegisterFunc registers a new function available for use within generated
// UI elements.
//
// The type of fn must be a function type; if it is not then the function will
// return an error. If a name is repeated, the function will be replaced in the
// set of available functions.
func (l *Loader) RegisterFunc(name string, fn interface{}) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return FunctionTypeError{
			Func: fn,
		}
	}
	l.callbacks[name] = fn
	return nil
}
