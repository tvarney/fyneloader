package fyneloader_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tvarney/fyneloader"
)

func TestReader(t *testing.T) {
	t.Parallel()
	t.Run("RegisterFunc", func(t *testing.T) {
		t.Parallel()
		t.Run("Function", func(t *testing.T) {
			t.Parallel()
			r := fyneloader.New()
			require.NoError(t, r.RegisterFunc("test", func() {}))
		})
		t.Run("NonFunction", func(t *testing.T) {
			t.Parallel()
			r := fyneloader.New()
			require.EqualError(
				t, r.RegisterFunc("test", true),
				fyneloader.FunctionTypeError{Func: true}.Error(),
			)
		})
	})
}
