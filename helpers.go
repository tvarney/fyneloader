package fyneloader

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
	"github.com/tvarney/maputil"
	"github.com/tvarney/maputil/errctx"
)

// GetFnVoidToVoid fetches a func() from the registered functions in the loader.
func GetFnVoidToVoid(l *Loader, data map[string]interface{}, key string) (func(), error) {
	fnname, ok, err := maputil.GetString(data, key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	fni, err := l.GetFunc(fnname)
	if err != nil {
		return nil, err
	}

	fn, ok := fni.(func())
	if !ok {
		return nil, FunctionTypeError{
			Func: fni,
		}
	}
	return fn, nil
}

// GetFnBoolToVoid fetches a func(bool) from the registered functions in the loader.
func GetFnBoolToVoid(l *Loader, data map[string]interface{}, key string) (func(bool), error) {
	fnname, ok, err := maputil.GetString(data, key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	fni, err := l.GetFunc(fnname)
	if err != nil {
		return nil, err
	}

	fn, ok := fni.(func(bool))
	if !ok {
		return nil, FunctionTypeError{
			Func: fni,
		}
	}
	return fn, nil
}

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

// GetTextStyle fetches and interprets a string from the map as a text style.
func GetTextStyle(data map[string]interface{}, key string) (fyne.TextStyle, error) {
	value, ok, err := maputil.GetString(data, key)
	if err != nil || !ok {
		return fyne.TextStyle{}, err
	}
	switch value {
	case "bold":
		return fyne.TextStyle{Bold: true, Italic: false, Monospace: false}, nil
	case "italic":
		return fyne.TextStyle{Bold: false, Italic: true, Monospace: false}, nil
	case "monospace":
		return fyne.TextStyle{Bold: false, Italic: false, Monospace: true}, nil
	case "bold+italic", "italic+bold":
		return fyne.TextStyle{Bold: true, Italic: true, Monospace: false}, nil
	}
	return fyne.TextStyle{}, maputil.EnumStringError{
		Value: value,
		Enum:  []string{"bold", "italic", "monospace", "bold+italic", "italic+bold"},
	}
}

func GetImage(ctx *errctx.Context, data map[string]interface{}) *canvas.Image {
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
