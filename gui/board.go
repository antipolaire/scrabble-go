package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"game"
	"go.uber.org/zap"
	"image/color"
)

type BoardWidget struct {
	fyne.Container
	fyne.Draggable
	Board               *game.Board
	IsDragging          bool
	CurrentDragPosition fyne.Position
	tileDragger         *TileDragger
	tilesByIndex        map[int]*game.Tile
	numColumns          int
	numRows             int
}

type BoardRenderer struct {
	BoardWidget *BoardWidget
	objects     []fyne.CanvasObject
}

func (b BoardRenderer) Destroy() {

}

func (b BoardRenderer) Layout(size fyne.Size) {

}

func (b BoardRenderer) MinSize() fyne.Size {
	return fyne.NewSize(CellWidth*16, CellHeight*17)
}

func (b BoardRenderer) Objects() []fyne.CanvasObject {
	return b.objects
}

func (b BoardRenderer) Refresh() {
}

func IsPositionInObjectBounds(canvasObject fyne.CanvasObject, position fyne.Position) bool {
	// TODO correctly implement this
	widgetPosition := canvasObject.Position()
	widgetSize := canvasObject.Size()
	if position.X < widgetPosition.X {
		return false
	}
	if position.Y < widgetPosition.Y {
		return false
	}
	if position.X > widgetPosition.X+widgetSize.Width {
		return false
	}
	if position.Y > widgetPosition.Y+widgetSize.Height {
		return false
	}
	zap.S().Info(fmt.Sprintf("Pos %v is in widget pos %v", position, widgetPosition))
	return true
}

func (b *BoardWidget) CreateRenderer() fyne.WidgetRenderer {
	return &BoardRenderer{
		BoardWidget: b,
		objects:     b.Container.Objects,
	}
}

func (b *BoardWidget) Dragged(event *fyne.DragEvent) {
	dragPosition := event.Position
	if !b.IsDragging {
		if b.tileDragger.OnDragStart(dragPosition) {
			b.IsDragging = true
		}
	} else {
		zap.L().Debug(fmt.Sprintf("Drag to %v", dragPosition))
		b.tileDragger.OnDrag(dragPosition)
	}
	b.CurrentDragPosition = dragPosition
}

func (b *BoardWidget) DragEnd() {
	if b.IsDragging {
		b.IsDragging = false
		b.tileDragger.OnDragEnd(b.CurrentDragPosition)
	}
	zap.S().Info("DragEnd")
}

func NewBoardWidget(myGame *game.Game) *BoardWidget {
	numBoardCols := 15
	numBoardRows := 15
	numCols := numBoardCols + NumIndexCols
	numRows := numBoardRows + NumIndexRows + NumRackRows
	cellStacks := make([]fyne.CanvasObject, numCols*numRows)
	tilesByIndex := make(map[int]*game.Tile, numCols*numRows)

	for i := 0; i < numCols; i++ {
		for j := 0; j < numRows; j++ {
			cellIndex := XY2I(i, j, numCols)
			// Rack row
			if j > numBoardRows {
				fieldColor := canvas.NewRectangle(color.White)
				fieldColor.StrokeColor = color.Black
				fieldColor.StrokeWidth = 1
				stack := container.NewStack(fieldColor)
				// Add tilesWidgets to stack for the first 7 tiles
				if i < 7 {
					tile := &myGame.CurrentPlayer.Tiles[i]
					tilesByIndex[cellIndex] = tile
					tileWidget := NewTileWidget(tile, myGame)
					stack.Add(tileWidget)
				}
				cellStacks[cellIndex] = stack
				continue
			}

			// Upper left corner cell
			if i == 0 && j == 0 {
				cornerRect := canvas.NewRectangle(color.White)
				cornerRect.StrokeColor = color.Black
				cornerRect.StrokeWidth = 1
				cellStacks[cellIndex] = cornerRect
				continue
			}

			// Column labels "A" to "O"
			if i == 0 && j > 0 {
				label := canvas.NewText(fmt.Sprintf("%c", 'A'+j-1), color.Black)
				label.Alignment = fyne.TextAlignCenter
				label.TextStyle = fyne.TextStyle{Bold: true}
				label.TextSize = 20
				label.Color = color.Black
				fieldColor := canvas.NewRectangle(color.White)
				fieldColor.StrokeColor = color.Black
				fieldColor.StrokeWidth = 1
				stack := container.NewStack(fieldColor, label)
				cellStacks[cellIndex] = stack
				continue
			}

			// Row labels "1" to "15"
			if j == 0 && i > 0 {
				label := canvas.NewText(fmt.Sprintf("%d", i), color.Black)
				label.Alignment = fyne.TextAlignCenter
				label.TextStyle = fyne.TextStyle{Bold: true}
				label.TextSize = 20
				label.Color = color.Black
				fieldColor := canvas.NewRectangle(color.White)
				fieldColor.StrokeColor = color.Black
				fieldColor.StrokeWidth = 1
				stack := container.NewStack(fieldColor, label)
				cellStacks[cellIndex] = stack
				continue
			}

			// Correct field index
			fieldX := i
			fieldY := j
			if i > 0 {
				fieldX--
			}
			if j > 0 {
				fieldY--
			}

			field, _ := myGame.Board.GetField(fieldX, fieldY)

			fieldColor, fieldText := CrateFieldStackComponents(field)
			stack := container.NewStack(fieldColor, fieldText)
			cellStacks[cellIndex] = stack
		}
	}

	boardWidget := &BoardWidget{
		Board:        myGame.Board,
		Container:    *container.New(layout.NewGridLayout(numCols), cellStacks...),
		numColumns:   numBoardCols,
		numRows:      numBoardRows,
		tilesByIndex: tilesByIndex,
		IsDragging:   false,
	}
	boardWidget.tileDragger = NewTileDragger(boardWidget, myGame)
	return boardWidget
}

func IntToRGBA(i int) color.RGBA {
	return color.RGBA{R: uint8(i >> 16), G: uint8(i >> 8), B: uint8(i), A: 0xff}
}

func FieldTypeToColor(field game.Field) color.Color {
	switch field.Type {
	case game.NF:
		return IntToRGBA(game.NFColor)
	case game.DL:
		return IntToRGBA(game.DLColor)
	case game.TL:
		return IntToRGBA(game.TLColor)
	case game.DW:
		return IntToRGBA(game.DWColor)
	case game.TW:
		return IntToRGBA(game.TWColor)
	case game.CS:
		return IntToRGBA(game.CSColor)
	}
	return color.White
}

func CrateFieldStackComponents(field game.Field) (*canvas.Rectangle, *canvas.Text) {
	cellColor := canvas.NewRectangle(FieldTypeToColor(field))
	cellColor.StrokeColor = color.Black
	cellColor.StrokeWidth = 1

	fieldText := canvas.NewText("", color.Black)
	switch field.Type {
	case game.NF:
		fieldText.Text = ""
	case game.DL:
		fieldText.Text = "2L"
	case game.TL:
		fieldText.Text = "3L"
	case game.DW:
		fieldText.Text = "2W"
	case game.TW:
		fieldText.Text = "3W"
	case game.CS:
		fieldText.Text = "*"
	}
	fieldText.Alignment = fyne.TextAlignCenter
	fieldText.TextStyle = fyne.TextStyle{Bold: true}
	fieldText.TextSize = 20
	fieldText.Color = color.White

	return cellColor, fieldText
}

// A tile drag works as follows:
// 1. Method Dragged is called with the current position of the mouse stored in the event
// 2. From the position, the cell on the grid is calculated and retrieved
// 3. Each cell consists of a stack of canvas objects
// 4. If the top object in the stack is an instance of TileWidget, the object is draggable
// 5. If the object is draggable, the position of the object is set to the position of the mouse
// 6. The object is dropped, it is checked if it is dropped above a cell
// 7. If it is dropped above a cell, the cell is retrieved and checked if it is empty
// 8. If the cell is empty, the tile is placed on top of the stack of the cell, removed from it's original position
//    and the stack is refreshed.
// 9. If the cell is not empty, the tile is placed back to its original position.
// 10. If the tile is dropped anywhere else, it is placed back to its original position.

// TileDragger To implement tile dragging logic, a separate struct (to store necessary information) and interface (to execute the required logic) are created.
type TileDragger struct {
	// The tile that is dragged
	Tile *game.Tile
	// TileWidget the widget that is dragged
	TileWidget *TileWidget
	// The position of the tile
	Position fyne.Position
	// The position of the tile before it was dragged
	PreviousPosition fyne.Position
	// The cell (stack) the tile is currently on
	CurrentCell *fyne.Container
	// The cell (stack) the tile was on before it was dragged
	PreviousCell *fyne.Container
	// The board the tile is on
	Board *BoardWidget
	// The game the tile is on
	Game *game.Game
}

type TileDraggerHandler interface {
	OnDragStart(position fyne.Position) bool
	OnDrag(position fyne.Position)
	OnDragEnd(position fyne.Position)

	// Check if the given position is within the area of the players rack
	IsInRackArea(position fyne.Position) bool
	// Check if the given position is within the area of the board
	IsInBoardArea(position fyne.Position) bool
	// IsCellEmpty checks if the cell at the given position is empty
	IsCellEmpty(position fyne.Position) bool

	// GetTileWidgetByPosition returns the tile widget at the given position
	GetTileWidgetByPosition(position fyne.Position) (*TileWidget, bool)
	// GetCellIndexByPosition returns the index of the cell at the given position
	GetCellIndexByPosition(position fyne.Position) (int, int)
	// GetCellByPosition returns the cell (stack) at the given position
	GetCellByPosition(position fyne.Position) (*fyne.Container, bool)
	// AddTileToCell adds the given tile to the cell at the given position
	AddTileToCell(tile *game.Tile, position fyne.Position)
	// RemoveTileFromCell removes the tile from the cell at the given position
	RemoveTileFromCell(position fyne.Position)
	// SwapTiles swaps the tiles at the given positions
	SwapTiles(position1 fyne.Position, position2 fyne.Position)
	// RefreshCell refreshes the cell at the given position
	RefreshCell(position fyne.Position)
	// RefreshTile refreshes the tile at the given position
	RefreshTile(position fyne.Position)
	// RefreshBoard refreshes the board
	RefreshBoard()
}

func NewTileDragger(board *BoardWidget, game *game.Game) *TileDragger {
	return &TileDragger{
		Board: board,
		Game:  game,
	}
}

func (d *TileDragger) GetTileByPosition(position fyne.Position) *game.Tile {
	x, y := d.GetCellIndexByPosition(position)
	index := XY2I(x, y, 16)
	tile := d.Board.tilesByIndex[index]
	zap.L().Debug(fmt.Sprintf("Got tile %v from pos %d,%d at index %d", tile, x, y, index))
	return tile
}
func (d *TileDragger) OnDragStart(position fyne.Position) bool {
	d.PreviousPosition = position
	var isValid = false
	d.PreviousCell, isValid = d.GetCellByPosition(position)
	if !isValid {
		return false
	}
	d.Tile = d.GetTileByPosition(position)
	if d.Tile == nil {
		return false
	}
	d.TileWidget, isValid = d.GetTileWidgetByPosition(position)
	if !isValid {
		return false
	}
	d.RemoveTileFromCell(position)
	d.Board.Refresh()
	return isValid
}

func (d *TileDragger) OnDrag(position fyne.Position) {
	d.Position = position
	d.TileWidget.Move(position)
	d.TileWidget.Refresh()
}

func (d *TileDragger) OnDragEnd(position fyne.Position) {
	zap.L().Debug(fmt.Sprintf("Drag end at %v", position))
	d.Position = position
	x, y := d.GetCellIndexByPosition(position)
	d.CurrentCell, _ = d.GetCellByIndex(x, y)
	if d.IsRackCell(x, y) {
		zap.L().Debug("Dropped in rack area")
		if d.IsCellEmpty(x, y) {
			zap.L().Debug("Dropped on empty cell")
			d.AddTileToCell(d.Tile, position)
			d.SwapTiles(position, d.PreviousPosition)
		} else {
			zap.L().Debug("Dropped on non empty cell")
			d.AddTileToCell(d.Tile, d.PreviousPosition)
		}
	} else if d.IsBoardCell(x, y) {
		zap.L().Debug("Dropped in board area")
		if d.IsCellEmpty(x, y) {
			zap.L().Debug("Dropped on empty cell")
			d.AddTileToCell(d.Tile, position)
			d.SwapTiles(position, d.PreviousPosition)
		} else {
			zap.L().Debug("Dropped on non empty cell")
			d.AddTileToCell(d.Tile, d.PreviousPosition)
		}
	} else {
		zap.L().Debug("Dropped outside legal area")
		d.AddTileToCell(d.Tile, d.PreviousPosition)
	}
	d.Board.Refresh()
	d.Tile = nil
	d.TileWidget = nil
	d.CurrentCell = nil
	d.PreviousCell = nil
}

func (d *TileDragger) IsRackCell(x int, y int) bool {
	return y == 16 && x >= 1 && x <= 15
}

func (d *TileDragger) IsBoardCell(x int, y int) bool {
	return y >= 1 && y <= 15 && x >= 1 && x <= 15
}

func (d *TileDragger) IsCellEmpty(x int, y int) bool {
	cell, ok := d.GetCellByIndex(x, y)
	if !ok {
		return false
	}
	// if any of it is TileWidget, cell is not empty
	for _, cell := range cell.Objects {
		_, ok := cell.(*TileWidget)
		if ok {
			return false
		}
	}
	return true
}

func (d *TileDragger) GetTileWidgetByPosition(position fyne.Position) (*TileWidget, bool) {
	cell, ok := d.GetCellByPosition(position)
	if !ok {
		return nil, false
	}
	if len(cell.Objects) == 0 {
		return nil, false
	}
	tileWidget, ok := cell.Objects[len(cell.Objects)-1].(*TileWidget)
	if !ok {
		return nil, false
	}
	return tileWidget, true
}

func (d *TileDragger) GetCellIndexByPosition(position fyne.Position) (int, int) {
	effectiveBoardSize := d.Board.Size()
	width := uint16((effectiveBoardSize.Width + theme.Padding()) / 16.0)
	height := uint16((effectiveBoardSize.Height + theme.Padding()) / 17.0)
	x, y := P2C(position, width, height)
	//zap.L().Debug(
	//	fmt.Sprintf(
	//		"Board Size: %v, "+
	//			"Padding: %v, "+
	//			"Inner Padding: %v, "+
	//			"Calc Width: %v, "+
	//			"Calc Height: %v, "+
	//			"Position: %v, "+
	//			"X: %v, "+
	//			"Y: %v",
	//		effectiveBoardSize,
	//		theme.Padding(),
	//		theme.InnerPadding(),
	//		width,
	//		height,
	//		position,
	//		x,
	//		y,
	//	))
	return x, y
}

func (d *TileDragger) GetCellByPosition(position fyne.Position) (*fyne.Container, bool) {
	x, y := d.GetCellIndexByPosition(position)
	zap.L().Debug(fmt.Sprintf("Cell index (%d,%d) from position %v", x, y, position))
	return d.GetCellByIndex(x, y)
}

func (d *TileDragger) GetCellByIndex(x int, y int) (*fyne.Container, bool) {
	if x < 0 || y < 0 || x > d.Board.numColumns+NumIndexCols || y > d.Board.numRows+NumIndexRows+NumRackRows {
		return nil, false
	}
	cell := d.Board.Container.Objects[XY2I(x, y, d.Board.numColumns+NumIndexCols)]
	cellTyped, ok := cell.(*fyne.Container)
	if !ok {
		return nil, false
	}
	return cellTyped, true
}

func (d *TileDragger) AddTileToCell(tile *game.Tile, position fyne.Position) {
	cell, _ := d.GetCellByPosition(position)
	tileWidget := NewTileWidget(tile, d.Game)
	cell.Add(tileWidget)
}

func (d *TileDragger) RemoveTileFromCell(position fyne.Position) {
	cell, isValid := d.GetCellByPosition(position)
	zap.L().Debug(fmt.Sprintf("Got cell %v (valid:%v) at position %v", cell, isValid, position))
	cell.Objects = cell.Objects[:1]
}

func (d *TileDragger) SwapTiles(p1 fyne.Position, p2 fyne.Position) {
	p1x, p1y := d.GetCellIndexByPosition(p1)
	p2x, p2y := d.GetCellIndexByPosition(p2)

	i1 := XY2I(p1x, p1y, d.Board.numColumns+NumIndexCols)
	i2 := XY2I(p2x, p2y, d.Board.numColumns+NumIndexCols)
	if i1 == i2 {
		return
	}
	t1 := d.Board.tilesByIndex[i1]
	t2 := d.Board.tilesByIndex[i2]
	zap.L().Debug(fmt.Sprintf("Swapped %v@%d <=> %v@%d", t1, i2, t2, i2))
	d.Board.tilesByIndex[i2] = t1
	d.Board.tilesByIndex[i1] = t2
}

func (d *TileDragger) RefreshCell(position fyne.Position) {
	cell, _ := d.GetCellByPosition(position)
	cell.Refresh()
}

func (d *TileDragger) RefreshTile(position fyne.Position) {
	tileWidget, _ := d.GetTileWidgetByPosition(position)
	tileWidget.Refresh()
}

func (d *TileDragger) RefreshBoard() {
	d.Board.Refresh()
}
