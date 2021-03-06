package fyneloader

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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
	return InvalidWidgetType(ctx, v)
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
	return InvalidWidgetType(ctx, v)
}

// CreateCard creates a new Card widget.
func CreateCard(ctx *errctx.Context, l *Loader, v interface{}) fyne.CanvasObject {
	switch w := v.(type) {
	case string:
		return widget.NewCard("", "", nil)
	case map[string]interface{}:
		title, _, err := maputil.GetString(w, KeyTitle)
		ctx.ErrorWithKey(err, KeyTitle)

		subtitle, _, err := maputil.GetString(w, KeySubTitle)
		ctx.ErrorWithKey(err, KeySubTitle)

		var child fyne.CanvasObject
		childdata, _, err := maputil.GetObject(w, KeyChild)
		ctx.ErrorWithKey(err, KeyChild)
		if childdata != nil {
			ctx.Path.Add(mpath.Key(KeyChild))
			child = l.Unpack(ctx, childdata)
			ctx.Path.Pop()
		}

		hidden, _, err := maputil.GetBoolean(w, KeyHidden)
		ctx.ErrorWithKey(err, KeyHidden)

		card := widget.NewCard(title, subtitle, child)
		card.Hidden = hidden
		card.Image = GetImage(ctx, w)
		return card
	}
	return InvalidWidgetType(ctx, v)
}

// CreateCheck creates a new Check widget.
func CreateCheck(ctx *errctx.Context, l *Loader, v interface{}) fyne.CanvasObject {
	switch w := v.(type) {
	case string:
		return widget.NewCheck("", nil)
	case map[string]interface{}:
		text, _, err := maputil.GetString(w, KeyText)
		ctx.ErrorWithKey(err, KeyText)

		fn, err := GetFnBoolToVoid(l, w, KeyFunc)
		ctx.ErrorWithKey(err, KeyFunc)

		hidden, _, err := maputil.GetBoolean(w, KeyHidden)
		ctx.ErrorWithKey(err, KeyHidden)

		chk := widget.NewCheck(text, fn)
		chk.Hidden = hidden
		return chk
	}
	return InvalidWidgetType(ctx, v)
}

// CreateHBox creates a new HBox container.
func CreateHBox(ctx *errctx.Context, l *Loader, v interface{}) fyne.CanvasObject {
	return createBox(ctx, l, v, container.NewHBox)
}

// CreateHSpacer creates a new horizontal spacer.
func CreateHSpacer(ctx *errctx.Context, l *Loader, v interface{}) fyne.CanvasObject {
	return createSpacer(false, true)
}

// CreateLabel creates a new Label.
func CreateLabel(ctx *errctx.Context, l *Loader, v interface{}) fyne.CanvasObject {
	switch w := v.(type) {
	case string:
		return widget.NewLabel("")
	case map[string]interface{}:
		text, _, err := maputil.GetString(w, KeyText)
		ctx.ErrorWithKey(err, KeyText)

		align, err := GetStringEnumAsInt(
			w, KeyAlign, []string{ValueDefault, ValueLeading, ValueCenter, ValueTrailing},
			[]int{
				int(fyne.TextAlignLeading), int(fyne.TextAlignLeading), int(fyne.TextAlignCenter),
				int(fyne.TextAlignTrailing),
			}, int(fyne.TextAlignLeading),
		)
		ctx.ErrorWithKey(err, KeyAlign)

		wrap, err := GetStringEnumAsInt(
			w, KeyWrap, []string{ValueDefault, ValueOff, ValueTruncate, ValueBreak, ValueWord},
			[]int{
				int(fyne.TextWrapOff), int(fyne.TextWrapOff), int(fyne.TextTruncate), int(fyne.TextWrapBreak),
				int(fyne.TextWrapWord),
			}, int(fyne.TextWrapOff),
		)
		ctx.ErrorWithKey(err, KeyWrap)

		style, err := GetTextStyle(w, KeyStyle)
		ctx.ErrorWithKey(err, KeyStyle)

		lbl := widget.NewLabel(text)
		lbl.Alignment = fyne.TextAlign(align)
		lbl.Wrapping = fyne.TextWrap(wrap)
		lbl.TextStyle = style
		return lbl
	}
	return InvalidWidgetType(ctx, v)
}

// CreateRadioGroup creates a new RadioGroup widget.
func CreateRadioGroup(ctx *errctx.Context, l *Loader, v interface{}) fyne.CanvasObject {
	switch w := v.(type) {
	case string:
		return widget.NewRadioGroup(nil, nil)
	case map[string]interface{}:
		var opts []string
		iopts, _, err := maputil.GetArray(w, KeyOptions)
		ctx.ErrorWithKey(err, KeyOptions)
		if len(iopts) > 0 {
			opts = make([]string, 0, len(iopts))
			for _, v := range iopts {
				strval, ok := v.(string)
				if !ok {
					strval = fmt.Sprintf("%v", v)
				}
				opts = append(opts, strval)
			}
		}

		fn, err := GetFnStringToVoid(l, w, KeyFunc)
		ctx.ErrorWithKey(err, KeyFunc)

		hidden, _, err := maputil.GetBoolean(w, KeyHidden)
		ctx.ErrorWithKey(err, KeyHidden)

		disabled, _, err := maputil.GetBoolean(w, KeyDisabled)
		ctx.ErrorWithKey(err, KeyDisabled)

		required, _, err := maputil.GetBoolean(w, KeyRequired)
		ctx.ErrorWithKey(err, KeyRequired)

		selected, err := GetStringFromArray(w, KeySelected, opts)
		ctx.ErrorWithKey(err, KeySelected)

		orientation, err := GetStringEnumAsInt(
			w, KeyOrientation,
			[]string{ValueDefault, ValueHorizontal, ValueVertical},
			[]int{1, 0, 1}, 0,
		)
		ctx.ErrorWithKey(err, KeyOrientation)

		rgroup := widget.NewRadioGroup(opts, fn)
		rgroup.Required = required
		rgroup.Hidden = hidden
		rgroup.Selected = selected
		if disabled {
			rgroup.Disable()
		}
		if orientation == 0 {
			rgroup.Horizontal = true
		}
		return rgroup
	}
	return InvalidWidgetType(ctx, v)
}

// CreateSpacer creates a new spacer which expands both vertically and
// horizontally.
func CreateSpacer(ctx *errctx.Context, l *Loader, v interface{}) fyne.CanvasObject {
	return createSpacer(true, true)
}

// CreateSlider creates a new Slider widget.
func CreateSlider(ctx *errctx.Context, l *Loader, v interface{}) fyne.CanvasObject {
	switch w := v.(type) {
	case string:
		return widget.NewSlider(0.0, 100.0)
	case map[string]interface{}:
		min, _, err := maputil.GetNumber(w, KeyMin)
		ctx.ErrorWithKey(err, KeyMin)

		max, ok, err := maputil.GetNumber(w, KeyMax)
		ctx.ErrorWithKey(err, KeyMax)
		if !ok || err != nil {
			max = 100.0
		}

		step, ok, err := maputil.GetNumber(w, KeyStep)
		ctx.ErrorWithKey(err, KeyStep)
		if !ok || err != nil {
			step = 1.0
		}

		fn, err := GetFnFloat64ToVoid(l, w, KeyFunc)
		ctx.ErrorWithKey(err, KeyFunc)

		orientation, err := GetStringEnumAsInt(
			w, KeyOrientation,
			[]string{ValueDefault, ValueHorizontal, ValueVertical},
			[]int{int(widget.Horizontal), int(widget.Horizontal), int(widget.Vertical)},
			int(widget.Horizontal),
		)
		ctx.ErrorWithKey(err, KeyOrientation)

		hidden, _, err := maputil.GetBoolean(w, KeyHidden)
		ctx.ErrorWithKey(err, KeyHidden)

		slide := widget.NewSlider(min, max)
		slide.Step = step
		slide.Hidden = hidden
		slide.Orientation = widget.Orientation(orientation)
		slide.OnChanged = fn
		return slide
	}
	return InvalidWidgetType(ctx, v)
}

// CreateVBox creates a new HBox container.
func CreateVBox(ctx *errctx.Context, l *Loader, v interface{}) fyne.CanvasObject {
	return createBox(ctx, l, v, container.NewVBox)
}

// CreateVSpacer creates a new vertical spacer.
func CreateVSpacer(ctx *errctx.Context, l *Loader, v interface{}) fyne.CanvasObject {
	return createSpacer(true, false)
}

func createBox(
	ctx *errctx.Context, l *Loader, v interface{},
	fn func(...fyne.CanvasObject) *fyne.Container,
) fyne.CanvasObject {
	switch w := v.(type) {
	case string:
		return fn()
	case map[string]interface{}:
		cdata, _, err := maputil.GetArray(w, KeyChildren)
		ctx.ErrorWithKey(err, KeyChildren)
		if len(cdata) == 0 {
			return fn()
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

		box := fn(children...)
		box.Hidden = hidden
		return box
	}
	return InvalidWidgetType(ctx, v)
}

func createSpacer(vertical, horizontal bool) fyne.CanvasObject {
	s := &layout.Spacer{
		FixHorizontal: !horizontal,
		FixVertical:   !vertical,
	}
	return s
}
