package fyneloader

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/tvarney/maputil"
)

// CreateButton creates a new button using the data in v.
func CreateButton(l *Loader, v interface{}) (fyne.CanvasObject, error) {
	switch w := v.(type) {
	case string:
		return widget.NewButton("", nil), nil
	case map[string]interface{}:
		text, _, err := maputil.GetString(w, KeyText)
		if err != nil {
			return nil, err
		}
		fn, err := GetFnVoidToVoid(l, w, KeyFunc)
		if err != nil {
			return nil, err
		}
		return widget.NewButton(text, fn), nil
	}
	return nil, maputil.InvalidTypeError{
		Actual:   maputil.TypeName(v),
		Expected: []string{maputil.TypeString, maputil.TypeObject},
	}
}
