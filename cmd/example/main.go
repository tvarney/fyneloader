package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/tvarney/fyneloader"
	"github.com/tvarney/maputil/errctx"
)

func main() {
	os.Exit(NewApp(os.Args[0], os.Args[1:], os.Stdout, os.Stderr).Run())
}

type App struct {
	Name  string
	OutFp io.Writer
	ErrFp io.Writer

	files       []string
	roots       map[string]map[string]fyne.CanvasObject
	exit        int
	app         fyne.App
	window      fyne.Window
	loader      *fyneloader.Loader
	ctx         *errctx.Context
	currentfile string
	currentroot string
	viewsmenu   *fyne.Menu
}

func NewApp(name string, args []string, outfp, errfp io.Writer) *App {
	if outfp == nil {
		outfp = ioutil.Discard
	}
	if errfp == nil {
		errfp = ioutil.Discard
	}
	return &App{
		Name:  name,
		OutFp: outfp,
		ErrFp: errfp,

		files:     args,
		roots:     map[string]map[string]fyne.CanvasObject{},
		loader:    fyneloader.New(),
		ctx:       errctx.New(&errctx.ErrorPrinter{Stream: errfp}),
		viewsmenu: fyne.NewMenu("Views"),
	}
}

func (a *App) load() error {
	for _, f := range a.files {
		fmt.Fprintf(a.OutFp, "Loading %s\n", f)
		roots, err := a.loader.ReadFile(a.ctx, f)
		if err != nil {
			return err
		}
		a.roots[f] = roots

		if a.currentroot != "" {
			continue
		}

		var first string
		for k := range roots {
			if k < first || first == "" {
				first = k
			}
		}
		if first != "" {
			a.currentfile = f
			a.currentroot = first
		}
	}

	if len(a.files) == 1 {
		// Just create the single top-level
		a.viewsmenu.Items = a.makeviewmenu(a.files[0], a.roots[a.files[0]])
	} else {
		// Else create a set of menu items for the files
		views := make([]*fyne.MenuItem, 0, len(a.roots))
		for f, r := range a.roots {
			items := a.makeviewmenu(f, r)
			views = append(views, &fyne.MenuItem{
				Label:       f,
				IsSeparator: false,
				Action:      nil,
				ChildMenu:   fyne.NewMenu("", items...),
			})
		}
		a.viewsmenu.Items = views
	}

	if a.currentfile != "" {
		fmt.Fprintf(a.OutFp, "Setting current content to %s %s\n", a.currentfile, a.currentroot)
		first := a.roots[a.currentfile][a.currentroot]
		if first != nil {
			a.window.SetContent(first)
		}
	}

	return nil
}

func (a *App) newviewitem(file, root string) *fyne.MenuItem {
	return fyne.NewMenuItem(root, func() {
		if a.currentfile == file && a.currentroot == root {
			fmt.Fprintf(a.OutFp, "Already on %s: %s\n", file, root)
			return
		}

		a.currentfile = file
		a.currentroot = root
		content := a.roots[file][root]
		content.Show()
		a.window.SetContent(content)
	})
}

func (a *App) makeviewmenu(filename string, roots map[string]fyne.CanvasObject) []*fyne.MenuItem {
	items := make([]*fyne.MenuItem, 0, len(roots))
	for k := range roots {
		items = append(items, a.newviewitem(filename, k))
	}
	return items
}

func (a *App) Run() int {
	if len(a.files) == 0 {
		return 0
	}

	a.exit = 0
	a.app = app.New()
	a.window = a.app.NewWindow("FyneLoader Example")
	reload := desktop.CustomShortcut{KeyName: fyne.KeyR, Modifier: desktop.ControlModifier}
	a.window.Canvas().AddShortcut(&reload, func(shortcut fyne.Shortcut) {
		err := a.load()
		if err != nil {
			fmt.Fprintf(a.ErrFp, "Error: %v\n", err)
		}
	})

	menu := fyne.NewMainMenu(
		&fyne.Menu{
			Label: "File",
			Items: []*fyne.MenuItem{
				{Label: "Reload", Action: func() {
					err := a.load()
					if err != nil {
						fmt.Fprintf(a.ErrFp, "Error: %v\n", err)
					}
				}},
				{Label: "Quit", Action: func() {
					a.window.Close()
				}},
			},
		},
		a.viewsmenu,
	)
	a.window.SetMainMenu(menu)

	err := a.load()
	if err != nil {
		fmt.Fprintf(a.ErrFp, "Error: %v\n", err)
		return 1
	}

	a.window.ShowAndRun()
	return 0
}
