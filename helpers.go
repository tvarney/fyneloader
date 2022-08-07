package fyneloader

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/tvarney/maputil"
	"github.com/tvarney/maputil/errctx"
	"github.com/tvarney/maputil/mpath"
	"github.com/tvarney/maputil/unpack"
)

// GetChild
func (l *Loader) GetChild(ctx *errctx.Context, data map[string]interface{}) fyne.CanvasObject {
	raw, ok := data[KeyChild]
	if !ok || raw == nil {
		return nil
	}

	ctx.Path.Add(mpath.Key(KeyChild))
	child := l.Unpack(ctx, raw)
	ctx.Path.Pop()
	return child
}

// GetImage fetches and loads an image from a series of keys.
func (l *Loader) GetImage(ctx *errctx.Context, data map[string]interface{}) *canvas.Image {
	imgpath, pathok, err := maputil.GetString(data, KeyImagePath)
	ctx.ErrorWithKey(err, KeyImagePath)

	imguri, uriok, err := maputil.GetString(data, KeyImageURI)
	ctx.ErrorWithKey(err, KeyImageURI)

	imgfill, err := GetStringEnumAsInt(
		data, KeyImageFill,
		[]string{ValueDefault, ValueStretch, ValueContain, ValueOriginal},
		[]int{
			int(canvas.ImageFillStretch), int(canvas.ImageFillStretch),
			int(canvas.ImageFillContain), int(canvas.ImageFillOriginal),
		}, int(canvas.ImageFillStretch),
	)
	ctx.ErrorWithKey(err, KeyImageFill)

	var img *canvas.Image
	if pathok {
		if uriok {
			ctx.Error(ConflictingKeysError{Keys: []string{KeyImagePath, KeyImageURI}})
		}
		img = canvas.NewImageFromFile(imgpath)
	} else if uriok {
		if !l.FetchURIs {
			ctx.ErrorWithKey(ErrFetchURIDisabled, KeyImageURI)
			return nil
		}

		path, err := storage.ParseURI(imguri)
		if err == nil {
			img = canvas.NewImageFromURI(path)
		} else {
			ctx.ErrorWithKey(err, KeyImageURI)
		}
	}
	if img != nil {
		img.FillMode = canvas.ImageFill(imgfill)
	}
	return img
}

// GetFnBoolToVoid fetches a func(bool) from the registered functions in the
// loader.
func GetFnBoolToVoid(l *Loader, data map[string]interface{}, key string) (func(bool), error) {
	fnname, ok, err := maputil.GetString(data, key)
	if err != nil || !ok {
		return nil, err
	}
	fni, err := l.GetFunc(fnname)
	if err != nil {
		return nil, err
	}

	fn, ok := fni.(func(bool))
	if !ok {
		return nil, FunctionTypeError{Func: fni}
	}
	return fn, nil
}

// GetFnFloat64ToVoid fetches a func(float64) from the registered functions in
// the loader.
func GetFnFloat64ToVoid(l *Loader, data map[string]interface{}, key string) (func(float64), error) {
	fnname, ok, err := maputil.GetString(data, key)
	if err != nil || !ok {
		return nil, err
	}
	fni, err := l.GetFunc(fnname)
	if err != nil {
		return nil, err
	}

	fn, ok := fni.(func(float64))
	if !ok {
		return nil, FunctionTypeError{Func: fni}
	}
	return fn, nil
}

// GetFnStringToVoid fetches a func(string) from the registered functions in the
// loader.
func GetFnStringToVoid(l *Loader, data map[string]interface{}, key string) (func(string), error) {
	fnname, ok, err := maputil.GetString(data, key)
	if err != nil || !ok {
		return nil, err
	}
	fni, err := l.GetFunc(fnname)
	if err != nil {
		return nil, err
	}

	fn, ok := fni.(func(string))
	if !ok {
		return nil, FunctionTypeError{Func: fni}
	}
	return fn, nil
}

// GetFnVoidToVoid fetches a func() from the registered functions in the loader.
func GetFnVoidToVoid(l *Loader, data map[string]interface{}, key string) (func(), error) {
	fnname, ok, err := maputil.GetString(data, key)
	if err != nil || !ok {
		return nil, err
	}
	fni, err := l.GetFunc(fnname)
	if err != nil {
		return nil, err
	}

	fn, ok := fni.(func())
	if !ok {
		return nil, FunctionTypeError{Func: fni}
	}
	return fn, nil
}

// GetStringEnumAsInt fetches a string value from the map and converts it to an
// integer.
func GetStringEnumAsInt(data map[string]interface{}, key string, allowed []string, values []int, def int) (int, error) {
	value, ok, err := maputil.GetString(data, key)
	if !ok || err != nil {
		return def, err
	}
	for i, v := range allowed {
		if v == value {
			return values[i], nil
		}
	}
	return def, maputil.EnumStringError{Value: value, Enum: allowed}
}

// GetStringFromArray fetches either a string value or an integer value.
//
// If the value is a string, it must be present in the given array.
func GetStringFromArray(data map[string]interface{}, key string, opts []string) (string, error) {
	item, ok, err := maputil.GetString(data, key)
	if err != nil && ok {
		idx, _, err := maputil.GetInteger(data, KeySelected)
		if err != nil {
			return "", maputil.InvalidTypeError{
				Actual:   maputil.TypeName(data[key]),
				Expected: []string{maputil.TypeString, maputil.TypeInteger},
			}
		}
		original := idx
		if idx < 0 {
			idx = int64(len(opts)) + idx
		}
		if idx < 0 || idx > int64(len(opts)) {
			return "", ArrayIndexOutOfBoundsError{Index: original}
		}
		return opts[idx], nil
	}
	for _, v := range opts {
		if item == v {
			return item, nil
		}
	}
	return "", ErrInvalidOption
}

// GetTextStyle fetches and interprets a string from the map as a text style.
func GetTextStyle(ctx *errctx.Context, data map[string]interface{}, key string) fyne.TextStyle {
	value := unpack.OptionalStringEnum(ctx, data, key, []string{ValueDefault, "bold", "italic", "monospace", "bold+italic", "italic+bold"}, ValueDefault)
	switch value {
	default:
		return fyne.TextStyle{Bold: false, Italic: false, Monospace: false}
	case ValueDefault:
		return fyne.TextStyle{Bold: false, Italic: false, Monospace: false}
	case "bold":
		return fyne.TextStyle{Bold: true, Italic: false, Monospace: false}
	case "italic":
		return fyne.TextStyle{Bold: false, Italic: true, Monospace: false}
	case "monospace":
		return fyne.TextStyle{Bold: false, Italic: false, Monospace: true}
	case "bold+italic", "italic+bold":
		return fyne.TextStyle{Bold: true, Italic: true, Monospace: false}
	}
}

// GetButtonAlign
func GetButtonAlign(ctx *errctx.Context, data map[string]interface{}) widget.ButtonAlign {
	value := unpack.OptionalStringEnum(ctx, data, KeyAlign, []string{ValueDefault, ValueCenter, ValueLeading, ValueTrailing}, ValueDefault)
	switch value {
	default:
		return widget.ButtonAlignCenter
	case ValueDefault, ValueCenter:
		return widget.ButtonAlignCenter
	case ValueLeading:
		return widget.ButtonAlignLeading
	case ValueTrailing:
		return widget.ButtonAlignTrailing
	}
}

// GetButtonIconPlacement
func GetButtonIconPlacement(ctx *errctx.Context, data map[string]interface{}) widget.ButtonIconPlacement {
	value := unpack.OptionalStringEnum(ctx, data, KeyIconPlace, []string{ValueDefault, ValueLeading, ValueTrailing}, ValueDefault)
	switch value {
	default:
		return widget.ButtonIconLeadingText
	case ValueDefault, ValueLeading:
		return widget.ButtonIconLeadingText
	case ValueTrailing:
		return widget.ButtonIconTrailingText
	}
}

// GetButtonImportance
func GetButtonImportance(ctx *errctx.Context, data map[string]interface{}) widget.ButtonImportance {
	value := unpack.OptionalStringEnum(ctx, data, KeyImportance, []string{ValueDefault, ValueLow, ValueMedium, ValueHigh}, ValueDefault)
	switch value {
	default:
		return widget.MediumImportance
	case ValueDefault, ValueMedium:
		return widget.MediumImportance
	case ValueLow:
		return widget.LowImportance
	case ValueHigh:
		return widget.HighImportance
	}
}

func GetTextAlign(ctx *errctx.Context, data map[string]interface{}) fyne.TextAlign {
	value := unpack.OptionalStringEnum(ctx, data, KeyAlign, []string{ValueDefault, ValueLeading, ValueCenter, ValueTrailing}, ValueDefault)
	switch value {
	default:
		return fyne.TextAlignLeading
	case ValueDefault, ValueLeading:
		return fyne.TextAlignLeading
	case ValueCenter:
		return fyne.TextAlignCenter
	case ValueTrailing:
		return fyne.TextAlignTrailing
	}
}

func GetTextWrap(ctx *errctx.Context, data map[string]interface{}) fyne.TextWrap {
	value := unpack.OptionalStringEnum(ctx, data, KeyWrap, []string{ValueDefault, ValueOff, ValueTruncate, ValueBreak, ValueWord}, ValueDefault)
	switch value {
	default:
		return fyne.TextWrapOff
	case ValueDefault, ValueOff:
		return fyne.TextWrapOff
	case ValueTruncate:
		return fyne.TextTruncate
	case ValueBreak:
		return fyne.TextWrapBreak
	case ValueWord:
		return fyne.TextWrapWord
	}
}

func GetOrientation(ctx *errctx.Context, data map[string]interface{}, defval widget.Orientation) widget.Orientation {
	value := unpack.OptionalStringEnum(ctx, data, KeyOrientation, []string{ValueDefault, ValueHorizontal, ValueVertical}, ValueDefault)
	switch value {
	default:
		return defval
	case ValueDefault:
		return defval
	case ValueHorizontal:
		return widget.Horizontal
	case ValueVertical:
		return widget.Vertical
	}
}

func InvalidWidgetType(ctx *errctx.Context, v interface{}) fyne.CanvasObject {
	ctx.Error(maputil.InvalidTypeError{
		Actual:   maputil.TypeName(v),
		Expected: []string{maputil.TypeString, maputil.TypeObject},
	})
	return nil
}
