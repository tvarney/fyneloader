package fyneloader

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/tvarney/maputil"
	"github.com/tvarney/maputil/errctx"
	"github.com/tvarney/maputil/mpath"
)

// CreateAccordion creates a new Accordion widget.
func CreateAccordion(ctx *errctx.Context, l *Loader, v interface{}) fyne.CanvasObject {
	switch w := v.(type) {
	case string:
		return widget.NewAccordion()
	case map[string]interface{}:
		var items []*widget.AccordionItem
		idata, _, err := maputil.GetArray(w, KeyItems)
		ctx.ErrorWithKey(err, KeyItems)

		multiopen, _, err := maputil.GetBoolean(w, KeyMultiOpen)
		ctx.ErrorWithKey(err, KeyMultiOpen)

		if idata != nil {
			items = make([]*widget.AccordionItem, 0, len(idata))
			for i, value := range idata {
				// Make sure it's an object
				data, err := maputil.AsObject(value)
				if err != nil {
					ctx.ErrorWithIndex(err, i)
					continue
				}

				ctx.Path.Add(mpath.Index(i))
				title, _, err := maputil.GetString(data, KeyTitle)
				ctx.ErrorWithKey(err, KeyTitle)

				open, _, err := maputil.GetBoolean(data, KeyOpen)
				ctx.ErrorWithKey(err, KeyOpen)

				details, _, err := maputil.GetObject(data, KeyChild)
				ctx.ErrorWithKey(err, KeyChild)

				ctx.Path.Add(mpath.Key(KeyChild))
				child := l.Unpack(ctx, details)
				item := widget.NewAccordionItem(title, child)
				item.Open = open
				items = append(items, item)

				ctx.Path.PopN(2)
			}
		}

		a := widget.NewAccordion(items...)
		a.MultiOpen = multiopen
		return a
	}
	ctx.Error(maputil.InvalidTypeError{
		Actual:   maputil.TypeName(v),
		Expected: []string{maputil.TypeString, maputil.TypeObject},
	})
	return nil
}

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

		disabled, _, err := maputil.GetBoolean(w, KeyDisabled)
		ctx.ErrorWithKey(err, KeyDisabled)

		hidden, _, err := maputil.GetBoolean(w, KeyHidden)
		ctx.ErrorWithKey(err, KeyHidden)

		vCenter := int(widget.ButtonAlignCenter)
		vLeading := int(widget.ButtonAlignLeading)
		vTrailing := int(widget.ButtonAlignTrailing)

		align, err := GetStringEnumAsInt(
			w, KeyAlign, []string{ValueDefault, ValueCenter, ValueLeading, ValueTrailing},
			[]int{vCenter, vCenter, vLeading, vTrailing}, vCenter,
		)
		ctx.ErrorWithKey(err, KeyAlign)

		vLeading = int(widget.ButtonIconLeadingText)
		vTrailing = int(widget.ButtonIconTrailingText)
		iconAlign, err := GetStringEnumAsInt(
			w, KeyIconPlace, []string{ValueDefault, ValueLeading, ValueTrailing},
			[]int{vLeading, vLeading, vTrailing}, vLeading,
		)
		ctx.ErrorWithKey(err, KeyIconPlace)

		vLow := int(widget.LowImportance)
		vMedium := int(widget.MediumImportance)
		vHigh := int(widget.HighImportance)
		importance, err := GetStringEnumAsInt(
			w, KeyImportance, []string{ValueDefault, ValueLow, ValueMedium, ValueHigh},
			[]int{vMedium, vLow, vMedium, vHigh}, vMedium,
		)
		ctx.ErrorWithKey(err, KeyImportance)

		btn := widget.NewButton(text, fn)
		btn.Alignment = widget.ButtonAlign(align)
		btn.IconPlacement = widget.ButtonIconPlacement(iconAlign)
		btn.Importance = widget.ButtonImportance(importance)
		btn.Hidden = hidden
		if disabled {
			btn.Disable()
		}
		return btn
	}
	ctx.Error(maputil.InvalidTypeError{
		Actual:   maputil.TypeName(v),
		Expected: []string{maputil.TypeString, maputil.TypeObject},
	})
	return nil
}

// CreateHBox creates a new HBox container.
func CreateHBox(ctx *errctx.Context, l *Loader, v interface{}) fyne.CanvasObject {
	return createBox(ctx, l, v, container.NewHBox)
}

// CreateVBox creates a new HBox container.
func CreateVBox(ctx *errctx.Context, l *Loader, v interface{}) fyne.CanvasObject {
	return createBox(ctx, l, v, container.NewVBox)
}

func createBox(
	ctx *errctx.Context, l *Loader, v interface{},
	fn func(...fyne.CanvasObject) *fyne.Container,
) fyne.CanvasObject {
	switch w := v.(type) {
	case string:
		return container.NewVBox()
	case map[string]interface{}:
		cdata, _, err := maputil.GetArray(w, KeyChildren)
		ctx.ErrorWithKey(err, KeyChildren)
		if len(cdata) == 0 {
			return container.NewVBox()
		}

		hidden, _, err := maputil.GetBoolean(w, KeyHidden)
		ctx.ErrorWithKey(err, KeyHidden)

		children := make([]fyne.CanvasObject, 0, len(cdata))
		for i, c := range cdata {
			ctx.Path.Add(mpath.Index(i))
			child := l.Unpack(ctx, c)
			ctx.Path.Pop()
			if child != nil {
				children = append(children, child)
			}
		}

		box := container.NewVBox(children...)
		box.Hidden = hidden
		return box
	}
	ctx.Error(maputil.InvalidTypeError{
		Actual:   maputil.TypeName(v),
		Expected: []string{maputil.TypeString, maputil.TypeObject},
	})
	return nil
}
