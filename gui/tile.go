package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"game"
	"image/color"
)

type TileWidget struct {
	widget.BaseWidget
	Tile          *game.Tile
	Objects       []fyne.CanvasObject
	gameReference *game.Game
}

type TileWidgetRenderer struct {
	TileWidget *TileWidget
	objects    []fyne.CanvasObject
}

func (b *TileWidget) MinSize() fyne.Size {
	return fyne.NewSize(CellWidth, CellHeight)
}

func (b *TileWidget) Move(newPosition fyne.Position) {
	b.BaseWidget.Move(newPosition)
	//for _, obj := range b.Objects {
	//	obj.Move(newPosition)
	//}
}

func (b *TileWidgetRenderer) Refresh() {
	b.Layout(b.TileWidget.MinSize())
	canvas.Refresh(b.TileWidget)
}

func (b *TileWidgetRenderer) Destroy() {

}

func (b *TileWidgetRenderer) Layout(size fyne.Size) {
	for _, child := range b.TileWidget.Objects {
		child.Resize(fyne.NewSize(CellWidth, CellHeight))
		child.Move(fyne.NewPos(0, 0))
	}
}

func (b *TileWidgetRenderer) MinSize() fyne.Size {
	return b.TileWidget.MinSize()
}

func (b *TileWidgetRenderer) Objects() []fyne.CanvasObject {
	return b.objects
}

func (b *TileWidget) CreateRenderer() fyne.WidgetRenderer {
	return &TileWidgetRenderer{
		TileWidget: b,
		objects:    b.Objects,
	}
}

func NewTileWidget(tile *game.Tile, myGame *game.Game) *TileWidget {
	tileText, scoreText := CreateTileStackComponents(tile)
	baseTile := canvas.NewImageFromFile("../assets/base_tile.svg")
	tileWidget := &TileWidget{
		Tile: tile,
		Objects: []fyne.CanvasObject{
			baseTile,
			tileText,
			scoreText,
		},
		gameReference: myGame,
	}
	tileWidget.ExtendBaseWidget(tileWidget)
	return tileWidget
}

func CreateTileStackComponents(tile *game.Tile) (*canvas.Text, *fyne.Container) {
	tileText := canvas.NewText(tile.Letter, color.Black)
	tileText.Alignment = fyne.TextAlignCenter

	tileText.TextSize = 36
	tileText.Color = color.Black

	letterScore := canvas.NewText(fmt.Sprintf("%d ", tile.LetterScore), color.Black)
	letterScore.Alignment = fyne.TextAlignTrailing
	letterScore.TextSize = 18
	letterScore.Color = color.Black

	letterScoreBox := container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		letterScore,
	)

	return tileText, letterScoreBox
}
