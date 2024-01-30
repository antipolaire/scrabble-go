package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"game"
	"go.uber.org/zap"
	"gui"
)

func main() {

	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any

	undo := zap.ReplaceGlobals(logger)
	defer undo()
	//config.WriteConfig(config.NewConfig())
	windowSize := fyne.NewSize(gui.CellWidth*16+100, gui.CellHeight*16+146)

	myApp := app.New()

	//driver := myApp.Driver()
	//isPositionInWidgetBoundsCallback := func(widget fyne.CanvasObject, position fyne.Position) bool {
	//	widgetPosition := driver.AbsolutePositionForObject(widget)
	//	widgetSize := widget.Size()
	//	if position.X < widgetPosition.X {
	//		return false
	//	}
	//	if position.Y < widgetPosition.Y {
	//		return false
	//	}
	//	if position.X > widgetPosition.X+widgetSize.Width {
	//		return false
	//	}
	//	if position.Y > widgetPosition.Y+widgetSize.Height {
	//		return false
	//	}
	//	zap.S().Info(fmt.Sprintf("Pos %v is in widget pos %v", position, widgetPosition))
	//	return true
	//}

	myWindow := myApp.NewWindow("Lets Play Scrabble!")

	myGame := game.NewGame()

	myGame.Players = append(myGame.Players, *game.NewPlayer("Player 1"))

	myGame.PullNewTilesFromBag(myGame.CurrentPlayer)

	mainGrid := gui.NewBoardWidget(myGame)

	playButton := widget.NewButton("Zug spielen!", func() {
		scoredPoints := myGame.PlayTemporaryMoves(myGame.CurrentPlayer)
		if scoredPoints > 0 {
			zap.S().Info(fmt.Sprintf("Player '%s' scored %d points", myGame.CurrentPlayer.Name, scoredPoints))
		}
	})

	passButton := widget.NewButton("Passen!", func() {
		zap.S().Info("Pass pressed")
	})

	remainingTilesLabel := widget.NewLabel(fmt.Sprintf("Verbleibende Steine: %d", len(myGame.Bag.Tiles)))
	yourPointsLabel := widget.NewLabel(fmt.Sprintf("Deine Punkte: %d", 0))
	yourNameLabel := widget.NewLabel(fmt.Sprintf("Dein Name: %s", myGame.CurrentPlayer.Name))

	actionButtons := container.NewVBox(playButton, passButton, yourNameLabel, yourPointsLabel, remainingTilesLabel)

	mainLayout := container.NewBorder(nil, nil, nil, actionButtons, mainGrid)

	myWindow.Resize(windowSize)
	myWindow.SetContent(mainLayout)
	myWindow.ShowAndRun()
}
