package main

import (
	"github.com/nsf/termbox-go"
	"github.com/mwindels/go-rogue/geom"
)

//Display modes for variableContents field of canvases.
const (
	displayNothing uint = iota
	displayMainMenu
	displayNewGame
	displayEnvironment
	displayPause
	totalDisplayModes
)

//Overlays by state.  For all intents and purposes, this is read-only.
var (
	stateOverlays = [totalStates]overlay{
		addToOverlay(initOverlay(), initCanvas(opaque, '▓', termbox.ColorDefault, termbox.ColorBlack, geom.InitRectangle(0.0, 0.0, 1.0, 1.0), displayMainMenu, addLabels(initCanvasConstants(), initLabel("GO ROGUE", termbox.ColorDefault, termbox.ColorBlack, geom.InitPoint(0.5, 0.1), xAlignCentre, yAlignCentre)))),
		addToOverlay(initOverlay(), initCanvas(opaque, '▓', termbox.ColorDefault, termbox.ColorBlack, geom.InitRectangle(0.0, 0.0, 1.0, 1.0), displayNewGame, addLabels(initCanvasConstants(), initLabel("Start a New Game...", termbox.ColorDefault, termbox.ColorBlack, geom.InitPoint(0.5, 0.15), xAlignCentre, yAlignCentre)))),
		addToOverlay(initOverlay(), initCanvas(opaque, '▓', termbox.ColorDefault, termbox.ColorBlack, geom.InitRectangle(0.0, 0.0, 1.0, 1.0), displayEnvironment, initCanvasConstants())),
		addToOverlay(initOverlay(), initCanvas(opaque, '▓', termbox.ColorDefault, termbox.ColorBlack, geom.InitRectangle(0.0, 0.0, 1.0, 1.0), displayEnvironment, initCanvasConstants()), initCanvas(opaque, '▓', termbox.ColorDefault, termbox.ColorBlack, geom.InitRectangle(0.3, 0.3, 0.4, 0.4), displayPause, addLabels(initCanvasConstants(), initLabel("Game Paused", termbox.ColorDefault, termbox.ColorBlack, geom.InitPoint(0.5, 0.2), xAlignCentre, yAlignCentre)))),
	}
)

//Data type used for the variableContents maps.
type variableContentsKey struct {
	state uint
	subStateIndex uint
	displayMode uint
}

//A variableContents map which stores selectors (lists of labels affected by state).
var (
	selectorMap = map[variableContentsKey]([]label){
		variableContentsKey{state: stateMainMenu, subStateIndex: stateMainMenuSelectorIndex, displayMode: displayMainMenu}: alignLabels(geom.InitRectangle(0.5, 0.0, 0.0, 1.0), xAlignCentre, yAlignCentre, initLocationlessLabel("New Game", termbox.ColorDefault, termbox.ColorBlack), initLocationlessLabel("Options", termbox.ColorDefault, termbox.ColorBlack), initLocationlessLabel("Exit", termbox.ColorDefault, termbox.ColorBlack)),
		variableContentsKey{state: statePausedGame, subStateIndex: statePausedGameSelectorIndex, displayMode: displayPause}: alignLabels(geom.InitRectangle(0.5, 0.2, 0.0, 0.8), xAlignCentre, yAlignCentre, initLocationlessLabel("Resume", termbox.ColorDefault, termbox.ColorBlack), initLocationlessLabel("Quit to Main Menu", termbox.ColorDefault, termbox.ColorBlack)),
	}
)

/* Ohter (potential) maps will go here... */

//Functions to be called when a canvas is drawn.  Indexed by displayMode.  Some entries are supposed to be modified while the program is running, like the entry for displayEnvironment.
var (
	displayModeFunctions = [totalDisplayModes](func(geom.Rectangle)){
		func(border geom.Rectangle){},
		func(border geom.Rectangle){},
		func(border geom.Rectangle){},
		func(border geom.Rectangle){},
		func(border geom.Rectangle){},
	}
)