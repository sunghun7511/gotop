package widgets

import (
	tui "github.com/gizak/termui/v3"
)

type Widget interface {
	Update()
	HandleSignal(event tui.Event)
	GetUI() tui.Drawable
}
