package gui

// Create GameBoardLayout

import "fyne.io/fyne/v2"

// Declare conformity with Layout interface
var _ fyne.Layout = (*GameBoardLayout)(nil)

type GameBoardLayout struct {
	absolutObjectPositions map[fyne.CanvasObject]fyne.Position
}

type GameBoardLayoutObject interface {
	fyne.Layout
	SetAbsolutObjectPosition(obj fyne.CanvasObject, position fyne.Position)
	GetAbsolutObjectPosition(obj fyne.CanvasObject) fyne.Position
}

func NewGameBoardLayout() GameBoardLayoutObject {
	return &GameBoardLayout{
		absolutObjectPositions: make(map[fyne.CanvasObject]fyne.Position),
	}
}

func (m *GameBoardLayout) SetAbsolutObjectPosition(obj fyne.CanvasObject, position fyne.Position) {
	m.absolutObjectPositions[obj] = position
}

func (m *GameBoardLayout) GetAbsolutObjectPosition(obj fyne.CanvasObject) fyne.Position {
	return m.absolutObjectPositions[obj]
}

func (m *GameBoardLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	// Place objects according to their orientation relative to the previous object:
	objectPosition := fyne.NewPos(0, 0)

	for _, child := range objects {
		if absolutPosition, ok := m.absolutObjectPositions[child]; ok {
			objectPosition = absolutPosition
		}
		child.Resize(size)
		child.Move(objectPosition)
	}
}

func (m *GameBoardLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		minSize = minSize.Max(child.MinSize())
	}
	return minSize
}
