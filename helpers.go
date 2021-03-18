package fyneloader

import (
	"github.com/tvarney/maputil"
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
