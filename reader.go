package fyneloader

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/tvarney/maputil"
	"gopkg.in/yaml.v2"
)

// CreateElementFn is the function type used for the element creation
// callbacks.
type CreateElementFn func(*Loader, interface{}) (fyne.CanvasObject, error)

// Loader allows for loading UI definitions at runtime.
type Loader struct {
	callbacks map[string]interface{}
	elements  map[string]CreateElementFn
}

// New returns a new Loader instance.
func New() *Loader {
	return &Loader{
		callbacks: map[string]interface{}{},
		elements: map[string]CreateElementFn{
			"button": CreateButton,
		},
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

// GetFunc returns the function with the given name.
//
// If a function with the given name was not registered, this function will
// return nil and an error indicating such.
func (l *Loader) GetFunc(name string) (interface{}, error) {
	fn, ok := l.callbacks[name]
	if !ok {
		return nil, UndefinedFunctionError{Name: name}
	}
	return fn, nil
}

// ReadFile reads a file as either YAML or JSON.
func (l *Loader) ReadFile(path string) (map[string]fyne.CanvasObject, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case "yaml", "yml":
		return l.ReadFileYAML(path)
	case "json":
		return l.ReadFileJSON(path)
	}
	return nil, fmt.Errorf("%w %q", ErrUnknownFileExt, ext)
}

// ReadFileYAML reads a file as a YAML definition file.
func (l *Loader) ReadFileYAML(path string) (map[string]fyne.CanvasObject, error) {
	in, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	m, err := l.ReadYAML(in)
	in.Close()
	return m, err
}

// ReadYAML takes a Reader and interprets it as YAML data.
func (l *Loader) ReadYAML(in io.Reader) (map[string]fyne.CanvasObject, error) {
	var generic map[string]interface{}
	err := yaml.NewDecoder(in).Decode(&generic)
	if err != nil {
		return nil, err
	}
	return l.Unmarshal(generic)
}

// ReadFileJSON reads a file as a JSON definition file.
func (l *Loader) ReadFileJSON(path string) (map[string]fyne.CanvasObject, error) {
	in, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	m, err := l.ReadJSON(in)
	in.Close()
	return m, err
}

// ReadJSON takes a Reader and interprets it as JSON data.
func (l *Loader) ReadJSON(in io.Reader) (map[string]fyne.CanvasObject, error) {
	var generic map[string]interface{}
	err := json.NewDecoder(in).Decode(&generic)
	if err != nil {
		return nil, err
	}
	return l.Unmarshal(generic)
}

// Unmarshal takes a YAML or JSON map and creates a map of widgets from it.
func (l *Loader) Unmarshal(definitions map[string]interface{}) (map[string]fyne.CanvasObject, error) {
	widgets := make(map[string]fyne.CanvasObject, len(definitions))
	for k, v := range definitions {
		w, err := l.Unpack(v)
		if err != nil {
			return nil, err
		}
		if w != nil {
			widgets[k] = w
		}
	}
	return widgets, nil
}

// Unpack handles loading a single widget.
func (l *Loader) Unpack(v interface{}) (fyne.CanvasObject, error) {
	switch w := v.(type) {
	case map[string]interface{}:
		typename, ok, err := maputil.GetString(w, KeyType)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrNoWidgetType
		}
		cb, ok := l.elements[typename]
		if !ok {
			return nil, UnknownElementType{TypeName: typename}
		}
		return cb(l, w)
	case string:
		cb, ok := l.elements[w]
		if !ok {
			return nil, UnknownElementType{TypeName: w}
		}
		return cb(l, w)
	}
	return nil, maputil.InvalidTypeError{
		Actual:   maputil.TypeName(v),
		Expected: []string{maputil.TypeObject, maputil.TypeString},
	}
}
