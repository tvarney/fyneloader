package fyneloader

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/tvarney/maputil"
	"github.com/tvarney/maputil/errctx"
)

// CreateButton creates a new button using the data in v.
func CreateButton(ctx *errctx.Context, l *Loader, v interface{}) fyne.CanvasObject {
	switch w := v.(type) {
	case string:
		return widget.NewButton("", nil)
	case map[string]interface{}:
		text, _, err := maputil.GetString(w, KeyText)
		ctx.ErrorWithKey(err, KeyText)

		fn, err := GetFnVoidToVoid(l, w, KeyFunc)
		ctx.ErrorWithKey(err, KeyFunc)

		return widget.NewButton(text, fn)
	}
	ctx.Error(maputil.InvalidTypeError{
		Actual:   maputil.TypeName(v),
		Expected: []string{maputil.TypeString, maputil.TypeObject},
	})
	return nil
}
