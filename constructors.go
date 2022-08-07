package fyneloader

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/tvarney/maputil"
	"github.com/tvarney/maputil/errctx"
	"github.com/tvarney/maputil/mpath"
	"github.com/tvarney/maputil/unpack"
)

// CreateAccordion creates a new Accordion widget.
func CreateAccordion(ctx *errctx.Context, l *Loader, data map[string]interface{}) fyne.CanvasObject {
	if data == nil {
		return widget.NewAccordion()
	}

	idata := unpack.OptionalArray(ctx, data, KeyItems, nil)
	var items []*widget.AccordionItem
	if idata != nil {
		items = make([]*widget.AccordionItem, 0, len(idata))
		for i, value := range idata {
			// Make sure it's an object
			raw, err := maputil.AsObject(value)
			if err != nil {
				ctx.ErrorWithIndex(err, i)
				continue
			}

			ctx.Path.Add(mpath.Index(i))

			item := widget.NewAccordionItem(unpack.OptionalString(ctx, data, KeyTitle, ""), l.GetChild(ctx, data))
			item.Open = unpack.OptionalBoolean(ctx, raw, KeyOpen, false)
			items = append(items, item)

			ctx.Path.Pop()
		}
	}

	a := widget.NewAccordion(items...)
	a.MultiOpen = unpack.OptionalBoolean(ctx, data, KeyMultiOpen, false)

	return a
}

// CreateButton creates a new button using the data in v.
func CreateButton(ctx *errctx.Context, l *Loader, data map[string]interface{}) fyne.CanvasObject {
	if data == nil {
		return widget.NewButton("", nil)
	}

	fn, err := GetFnVoidToVoid(l, data, KeyFunc)
	ctx.ErrorWithKey(err, KeyFunc)

	btn := widget.NewButton(unpack.OptionalString(ctx, data, KeyText, ""), fn)
	btn.Alignment = GetButtonAlign(ctx, data)
	btn.IconPlacement = GetButtonIconPlacement(ctx, data)
	btn.Importance = GetButtonImportance(ctx, data)
	btn.Hidden = unpack.OptionalBoolean(ctx, data, KeyHidden, false)
	if unpack.OptionalBoolean(ctx, data, KeyDisabled, false) {
		btn.Disable()
	}
	return btn
}

// CreateCard creates a new Card widget.
func CreateCard(ctx *errctx.Context, l *Loader, data map[string]interface{}) fyne.CanvasObject {
	if data == nil {
		return widget.NewCard("", "", nil)
	}

	card := widget.NewCard(
		unpack.OptionalString(ctx, data, KeyTitle, ""),
		unpack.OptionalString(ctx, data, KeySubTitle, ""),
		l.GetChild(ctx, data),
	)
	card.Hidden = unpack.OptionalBoolean(ctx, data, KeyHidden, false)
	card.Image = l.GetImage(ctx, data)
	return card
}

// CreateCheck creates a new Check widget.
func CreateCheck(ctx *errctx.Context, l *Loader, data map[string]interface{}) fyne.CanvasObject {
	if data == nil {
		return widget.NewCheck("", nil)
	}

	fn, err := GetFnBoolToVoid(l, data, KeyFunc)
	ctx.ErrorWithKey(err, KeyFunc)

	check := widget.NewCheck(unpack.OptionalString(ctx, data, KeyText, ""), fn)
	check.Hidden = unpack.OptionalBoolean(ctx, data, KeyHidden, false)
	if unpack.OptionalBoolean(ctx, data, KeyDisabled, false) {
		check.Disable()
	}
	return check
}

// CreateHBox creates a new HBox container.
func CreateHBox(ctx *errctx.Context, l *Loader, data map[string]interface{}) fyne.CanvasObject {
	return createBox(ctx, l, data, container.NewHBox)
}

// CreateHSpacer creates a new horizontal spacer.
func CreateHSpacer(ctx *errctx.Context, l *Loader, data map[string]interface{}) fyne.CanvasObject {
	return createSpacer(false, true)
}

// CreateLabel creates a new Label.
func CreateLabel(ctx *errctx.Context, l *Loader, data map[string]interface{}) fyne.CanvasObject {
	if data == nil {
		return widget.NewLabel("")
	}

	label := widget.NewLabel(unpack.OptionalString(ctx, data, KeyText, ""))
	label.Alignment = GetTextAlign(ctx, data)
	label.Wrapping = GetTextWrap(ctx, data)
	label.TextStyle = GetTextStyle(ctx, data, KeyStyle)
	return label
}

// CreateRadioGroup creates a new RadioGroup widget.
func CreateRadioGroup(ctx *errctx.Context, l *Loader, data map[string]interface{}) fyne.CanvasObject {
	if data == nil {
		return widget.NewRadioGroup(nil, nil)
	}

	fn, err := GetFnStringToVoid(l, data, KeyFunc)
	ctx.ErrorWithKey(err, KeyFunc)
	rgroup := widget.NewRadioGroup(unpack.OptionalStringArray(ctx, data, KeyOptions), fn)
	rgroup.Hidden = unpack.OptionalBoolean(ctx, data, KeyHidden, false)
	rgroup.Required = unpack.OptionalBoolean(ctx, data, KeyRequired, false)
	rgroup.Selected = unpack.OptionalString(ctx, data, KeySelected, "")
	rgroup.Horizontal = GetOrientation(ctx, data, widget.Vertical) == widget.Horizontal
	if unpack.OptionalBoolean(ctx, data, KeyDisabled, false) {
		rgroup.Disable()
	}
	return rgroup
}

// CreateSpacer creates a new spacer which expands both vertically and
// horizontally.
func CreateSpacer(ctx *errctx.Context, l *Loader, data map[string]interface{}) fyne.CanvasObject {
	return createSpacer(true, true)
}

// CreateSlider creates a new Slider widget.
func CreateSlider(ctx *errctx.Context, l *Loader, data map[string]interface{}) fyne.CanvasObject {
	if data == nil {
		return widget.NewSlider(0.0, 100.0)
	}

	fn, err := GetFnFloat64ToVoid(l, data, KeyFunc)
	ctx.ErrorWithKey(err, KeyFunc)

	slider := widget.NewSlider(
		unpack.OptionalNumber(ctx, data, KeyMin, 0.0),
		unpack.OptionalNumber(ctx, data, KeyMax, 100.0),
	)
	slider.Step = unpack.OptionalNumber(ctx, data, KeyStep, 1.0)
	slider.OnChanged = fn
	slider.Orientation = GetOrientation(ctx, data, widget.Horizontal)
	slider.Hidden = unpack.OptionalBoolean(ctx, data, KeyHidden, false)
	return slider
}

// CreateVBox creates a new HBox container.
func CreateVBox(ctx *errctx.Context, l *Loader, data map[string]interface{}) fyne.CanvasObject {
	return createBox(ctx, l, data, container.NewVBox)
}

// CreateVSpacer creates a new vertical spacer.
func CreateVSpacer(ctx *errctx.Context, l *Loader, data map[string]interface{}) fyne.CanvasObject {
	return createSpacer(true, false)
}

func createBox(
	ctx *errctx.Context, l *Loader, data map[string]interface{},
	fn func(...fyne.CanvasObject) *fyne.Container,
) fyne.CanvasObject {
	if data == nil {
		return fn()
	}
	raw := unpack.OptionalArray(ctx, data, KeyChildren, nil)
	var children []fyne.CanvasObject
	if len(raw) > 0 {
		children = make([]fyne.CanvasObject, 0, len(raw))
		for i, c := range raw {
			ctx.Path.Add(mpath.Index(i))
			child := l.Unpack(ctx, c)
			ctx.Path.Pop()
			if child != nil {
				children = append(children, child)
			}
		}
	}
	box := fn(children...)
	box.Hidden = unpack.OptionalBoolean(ctx, data, KeyHidden, false)
	return box
}

func createSpacer(vertical, horizontal bool) fyne.CanvasObject {
	s := &layout.Spacer{
		FixHorizontal: !horizontal,
		FixVertical:   !vertical,
	}
	return s
}
