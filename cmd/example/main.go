package main

import (
	"fmt"
	"io"
	"os"

	"fyne.io/fyne/v2/app"
	"github.com/tvarney/fyneloader"
	"github.com/tvarney/maputil/errctx"
)

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}

func run(args []string, outfp, errfp io.Writer) int {
	if len(args) == 1 {
		return 0
	}

	a := app.New()
	window := a.NewWindow("Example")

	loader := fyneloader.New()
	ctx := errctx.New(&errctx.ErrorPrinter{Stream: errfp})

	for _, f := range args[1:] {
		roots, err := loader.ReadFile(ctx, f)
		if err != nil {
			fmt.Fprintf(errfp, "Error: %v\n", err)
			return 1
		}

		for _, root := range roots {
			window.SetContent(root)
			window.Show()
			a.Run()
		}
	}

	return 0
}
