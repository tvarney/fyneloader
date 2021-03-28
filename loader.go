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
	"github.com/tvarney/maputil/errctx"
	"github.com/tvarney/maputil/mpath"
	"gopkg.in/yaml.v3"
)

// CreateElementFn is the function type used for the element creation
// callbacks.
type CreateElementFn func(*errctx.Context, *Loader, interface{}) fyne.CanvasObject

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
			"accordion": CreateAccordion,
			"button":    CreateButton,
			"card":      CreateCard,
			"check":     CreateCheck,
			"hbox":      CreateHBox,
			"hspacer":   CreateHSpacer,
			"label":     CreateLabel,
			"spacer":    CreateSpacer,
			"vbox":      CreateVBox,
			"vspacer":   CreateVSpacer,
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
func (l *Loader) ReadFile(ctx *errctx.Context, path string) (map[string]fyne.CanvasObject, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		return l.ReadFileYAML(ctx, path)
	case ".json":
		return l.ReadFileJSON(ctx, path)
	}
	return nil, fmt.Errorf("%w %q", ErrUnknownFileExt, ext)
}

// ReadFileYAML reads a file as a YAML definition file.
func (l *Loader) ReadFileYAML(ctx *errctx.Context, path string) (map[string]fyne.CanvasObject, error) {
	in, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	m, err := l.ReadYAML(ctx, in)
	in.Close()
	return m, err
}

// ReadYAML takes a Reader and interprets it as YAML data.
func (l *Loader) ReadYAML(ctx *errctx.Context, in io.Reader) (map[string]fyne.CanvasObject, error) {
	var generic map[string]interface{}
	err := yaml.NewDecoder(in).Decode(&generic)
	if err != nil {
		return nil, err
	}
	return l.Unmarshal(ctx, generic)
}

// ReadFileJSON reads a file as a JSON definition file.
func (l *Loader) ReadFileJSON(ctx *errctx.Context, path string) (map[string]fyne.CanvasObject, error) {
	in, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	m, err := l.ReadJSON(ctx, in)
	in.Close()
	return m, err
}

// ReadJSON takes a Reader and interprets it as JSON data.
func (l *Loader) ReadJSON(ctx *errctx.Context, in io.Reader) (map[string]fyne.CanvasObject, error) {
	var generic map[string]interface{}
	err := json.NewDecoder(in).Decode(&generic)
	if err != nil {
		return nil, err
	}
	return l.Unmarshal(ctx, generic)
}

// Unmarshal takes a YAML or JSON map and creates a map of widgets from it.
func (l *Loader) Unmarshal(ctx *errctx.Context, data map[string]interface{}) (map[string]fyne.CanvasObject, error) {
	if ctx == nil {
		// New empty context
		ctx = errctx.New()
	}
	ctx.Reset()
	widgets := make(map[string]fyne.CanvasObject, len(data))
	for k, v := range data {
		ctx.Path.Add(mpath.Key(k))
		w := l.Unpack(ctx, v)
		if w != nil {
			widgets[k] = w
		}
		ctx.Path.Pop()
	}
	return widgets, nil
}

// Unpack handles loading a single widget.
func (l *Loader) Unpack(ctx *errctx.Context, v interface{}) fyne.CanvasObject {
	switch w := v.(type) {
	case map[string]interface{}:
		typename, ok, err := maputil.GetString(w, KeyType)
		if err != nil {
			ctx.ErrorWithKey(err, KeyType)
			return nil
		}
		if !ok {
			ctx.Error(ErrNoWidgetType)
			return nil
		}

		cb, ok := l.elements[typename]
		if !ok {
			ctx.ErrorWithKey(UnknownElementType{TypeName: typename}, KeyType)
			return nil
		}
		return cb(ctx, l, w)
	case string:
		cb, ok := l.elements[w]
		if !ok {
			ctx.Error(UnknownElementType{TypeName: w})
			return nil
		}
		return cb(ctx, l, w)
	}
	ctx.Error(maputil.InvalidTypeError{
		Actual:   maputil.TypeName(v),
		Expected: []string{maputil.TypeObject, maputil.TypeString},
	})
	return nil
}
