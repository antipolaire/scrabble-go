package gui

import "fyne.io/fyne/v2"

// Event type used as callback for widgets to check if a given position (e.g. from a drag event) is within the absolut
// bounds of the widget
type IsPositionInWidgetBounds func(widget fyne.CanvasObject, position fyne.Position) bool
