package main

import (
	"time"
	"math"
	"github.com/nsf/termbox-go"
	"github.com/mwindels/go-rogue/geom"
)

//Sometimes the maxLineLen can be a little off due to rounding errors on canvases smaller than the screen size (maxLines is probably also affected)
func drawLabel(border geom.Rectangle, lbl label) {
	lblLen := float64(len([]rune(lbl.text)))
	lblPoint := geom.InitPoint(math.Min(math.Max(border.UpperLeft().X + lbl.location.X * (border.UpperRight().X - border.UpperLeft().X), math.Floor(border.UpperLeft().X) + 1), math.Floor(border.UpperRight().X) - 1),
							math.Min(math.Max(border.UpperLeft().Y + lbl.location.Y * (border.LowerLeft().Y - border.UpperLeft().Y), math.Floor(border.UpperLeft().Y) + 1), math.Floor(border.LowerLeft().Y) - 1))
	maxLineLen := math.Floor(border.UpperRight().X - lblPoint.X)
	initialX := lblPoint.X
	if lbl.xAlign == xAlignCentre {
		maxLineLen = math.Floor(2 * math.Min(border.UpperRight().X - lblPoint.X, lblPoint.X - border.UpperLeft().X)) - 1
		initialX = lblPoint.X - maxLineLen / 2 + 1
	}else if lbl.xAlign == xAlignRight {
		maxLineLen = math.Floor(lblPoint.X - border.UpperLeft().X)
		initialX = lblPoint.X - maxLineLen + 1
	}
	maxLines := math.Floor(border.LowerLeft().Y - lblPoint.Y)
	initialY := lblPoint.Y
	if lbl.yAlign == yAlignCentre {
		maxLines = math.Floor(2 * math.Min(border.LowerLeft().Y - lblPoint.Y, lblPoint.Y - border.UpperLeft().Y)) - 1
		initialY = lblPoint.Y - maxLines / 2 + 1
	}else if lbl.yAlign == yAlignAbove {
		maxLines = math.Floor(lblPoint.Y - border.UpperLeft().Y)
		initialY = lblPoint.Y - maxLines + 1
	}
	usedY := int(math.Min(math.Ceil(lblLen / maxLineLen), maxLines))
	for y := 0; y < usedY; y++ {
		usedX := int(maxLineLen)
		if float64(y + 1) * maxLineLen > lblLen {
			usedX = int(lblLen) % int(maxLineLen)
		}
		for x := 0; x < usedX; x++ {
			if usedX < int(maxLineLen) {
				termbox.SetCell(int(initialX) + int((float64(lbl.xAlign) / 2) * float64(int(maxLineLen) - usedX)) + x, int(initialY) + int((float64(lbl.yAlign) / 2) * float64(int(maxLines) - usedY)) + y, ([]rune(lbl.text))[y * int(maxLineLen) + x], lbl.textColor, lbl.textHighlight)
			}else{
				termbox.SetCell(int(initialX) + x, int(initialY) + int((float64(lbl.yAlign) / 2) * float64(int(maxLines) - usedY)) + y, ([]rune(lbl.text))[y * int(maxLineLen) + x], lbl.textColor, lbl.textHighlight)
			}
		}
	}
}

func drawCanvasConstants(border geom.Rectangle, cc canvasConstants) {
	for i := 0; i < len(cc.labels); i++ {
		drawLabel(border, cc.labels[i])
	}
}

func drawSelection(border geom.Rectangle, selections []label, selected uint) {
	for i := 0; i < len(selections); i++ {
		if i == int(selected) {
			emboldened := selections[i]
			emboldened.textColor ^= termbox.AttrBold
			drawLabel(border, emboldened)
		}else{
			drawLabel(border, selections[i])
		}
	}
}

func drawBorder(border geom.Rectangle, borderCell termbox.Cell) {
	for x := int(border.UpperLeft().X); x <= int(border.UpperRight().X); x++ {
		termbox.SetCell(x, int(border.UpperLeft().Y), borderCell.Ch, borderCell.Fg, borderCell.Bg)
		termbox.SetCell(x, int(border.LowerLeft().Y), borderCell.Ch, borderCell.Fg, borderCell.Bg)
	}
	for y := int(border.UpperLeft().Y); y <= int(border.LowerLeft().Y); y++ {
		termbox.SetCell(int(border.UpperLeft().X), y, borderCell.Ch, borderCell.Fg, borderCell.Bg)
		termbox.SetCell(int(border.UpperRight().X), y, borderCell.Ch, borderCell.Fg, borderCell.Bg)
	}
}

func drawOverlay(o overlay, state uint, stRcv <-chan stateDescriptor, stRqst chan<- stateRequest) {
	termWidth, termHeight := termbox.Size()
	termXScale, termYScale := float64(termWidth - 1), float64(termHeight - 1)
	for i := 0; i < len(o.canvases); i++ {
		
		//Draw Canvas Background
		scaledBorder := geom.ScaleRectangle(o.canvases[i].border, termXScale, termYScale)
		if (o.canvases[i].attributes & opaque) != 0 && canvasLayerOverlaps(o, i) {
			for x := int(scaledBorder.UpperLeft().X) + 1; x < int(scaledBorder.UpperRight().X); x++ {
				for y := int(scaledBorder.UpperLeft().Y) + 1; y < int(scaledBorder.LowerLeft().Y); y++ {
					termbox.SetCell(x, y, ' ', 0, 0)
				}
			}
		}
		
		//Draw Variable Canvas Contents
		func(){
			defer func(){recover()}()	//will only fire if state changes between when state was polled, and substate was polled
			for j := 0; j < int(totalSubStates[state]); j++ {
				if selectorMap[variableContentsKey{state: state, subStateIndex: uint(j), displayMode: o.canvases[i].variableContents}] != nil {
					drawSelection(scaledBorder, selectorMap[variableContentsKey{state: state, subStateIndex: uint(j), displayMode: o.canvases[i].variableContents}], uint(getSubState(stRqst, stRcv, state, uint(j))))
				}
			}
		}()
		displayModeFunctions[o.canvases[i].variableContents](scaledBorder)
		
		/*switch o.canvases[i].variableContents {
		case displayMainMenu:
			if state == stateMainMenu {
				func(){
					defer func(){recover()}()	//will only fire if state changes between when state was polled, and substate was polled
					subState := getSubState(stRqst, stRcv, state, stateMainMenuSelectorIndex)
					drawSelection(scaledBorder, YYY, uint(subState))		//need a place to store these labels (maybe store them as canvasConstants?)
				}()
			}
		case displayNewGame:
			
		case displayEnvironment:
			drawEnvironment(scaledBorder, env)
		case displayPause:
			if state == statePausedGame {
				func(){
					defer func(){recover()}()	//the comment above applies here, too
					subState := getSubState(stRqst, stRcv, state, statePausedGameSelectorIndex)
					drawSelection(scaledBorder, YYY2, uint(subState))
				}()
			}
		}*/
		
		//Draw Constant Canvas Contents
		drawCanvasConstants(scaledBorder, o.canvases[i].constantContents)
		
		//Draw Canvas Border
		if (o.canvases[i].attributes & borderless) == 0 {
			drawBorder(scaledBorder, o.canvases[i].borderCell)
		}
	}
}

func drawEnvironment(border geom.Rectangle, env *environment) {
	initialX := int(math.Max(math.Floor(border.UpperLeft().X + 1), math.Ceil(border.UpperLeft().X + (border.UpperRight().X - border.UpperLeft().X) / 2 - float64(env.width) / 2)))
	initialY := int(math.Max(math.Floor(border.UpperLeft().Y + 1), math.Ceil(border.UpperLeft().Y + (border.LowerLeft().Y - border.UpperLeft().Y) / 2 - float64(env.height) / 2)))
	env.mutex.RLock()
	for x := initialX; x < int(math.Min(border.UpperRight().X, float64(initialX + env.width))); x++ {
		for y := initialY; y < int(math.Min(border.LowerLeft().Y, float64(initialY + env.height))); y++ {
			if env.entities[x - initialX][y - initialY] != nil {
				termbox.SetCell(x, y, env.entities[x - initialX][y - initialY].symbol, env.entities[x - initialX][y - initialY].color, 0)		//Perhaps Bg should depend on the status of the entity?
			}else{
				termbox.SetCell(x, y, env.tiles[x - initialX][y - initialY].Ch, env.tiles[x - initialX][y - initialY].Fg, env.tiles[x - initialX][y - initialY].Bg)
			}
		}
	}
	env.mutex.RUnlock()
}

func runRenderer(envRcv <-chan *environment, envRqst chan<- bool, stRcv <-chan stateDescriptor, stRqst chan<- stateRequest) {
	envRqst <- true
	env := <- envRcv
	displayModeFunctions[displayEnvironment] = func(border geom.Rectangle){drawEnvironment(border, env)}
	
	for {
		time.Sleep(time.Second / 30)
		
		state := getState(stRqst, stRcv)
		drawOverlay(stateOverlays[state], state, stRcv, stRqst)
		
		err := termbox.Flush()
		if err != nil {
			panic(err)
		}
		err = termbox.Clear(0, 0)
		if err != nil {
			panic(err)
		}
	}
}