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
			l := fyneloader.New()
			require.NoError(t, l.RegisterFunc("test", func() {}))
		})
		t.Run("NonFunction", func(t *testing.T) {
			t.Parallel()
			l := fyneloader.New()
			require.EqualError(
				t, l.RegisterFunc("test", true),
				fyneloader.FunctionTypeError{Func: true}.Error(),
			)
		})
	})
}
