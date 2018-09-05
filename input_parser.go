package main

import (
	"math"
	"github.com/nsf/termbox-go"
	"github.com/mwindels/go-rogue/geom"
)

const (
	commandKey uint = iota
	commandMouse	//separate commands for move, attack, and keys?
)

type inputCommand struct {
	commandType uint
	key termbox.Key
	x int
	y int
}

func initKeyCommand(keyIn termbox.Key) inputCommand {
	return inputCommand{commandType: commandKey, key: keyIn, x: 0, y: 0}
}

func initMouseCommand(mouseIn termbox.Key, xIn, yIn int) inputCommand {
	return inputCommand{commandType: commandMouse, key: mouseIn, x: xIn, y: yIn}
}

//The calculations made in this function are all dependent on how the renderer draws borders and the environment.
func mapClickToEnvironment(click termbox.Event, env *environment, state uint) (bool, int, int) {
	termWidth, termHeight := termbox.Size()	//Possible that these may change between when the event was polled and when this call was made.
	for i := len(stateOverlays[int(state)].canvases) - 1; i >= 0; i-- {
		canvas := stateOverlays[int(state)].canvases[i]
		border := geom.ScaleRectangle(canvas.border, float64(termWidth - 1), float64(termHeight - 1))
		if int(border.UpperLeft().X) <= click.MouseX && click.MouseX <= int(border.UpperRight().X) && int(border.UpperLeft().Y) <= click.MouseY && click.MouseY <= int(border.LowerLeft().Y) {
			if int(border.UpperLeft().X) < click.MouseX && click.MouseX < int(border.UpperRight().X) && int(border.UpperLeft().Y) < click.MouseY && click.MouseY < int(border.LowerLeft().Y) || (canvas.attributes & borderless) != 0 {
				if canvas.variableContents == displayEnvironment {
					env.mutex.RLock()
					xLoc := env.player.x + (click.MouseX - int(math.Ceil(border.UpperLeft().X + border.UpperRight().X / 2)))
					yLoc := env.player.y + (click.MouseY - int(math.Ceil(border.UpperLeft().Y + border.LowerLeft().Y / 2)))
					env.mutex.RUnlock()
					return true, xLoc, yLoc
				}
			}
			if (canvas.attributes & opaque) != 0 {
				break
			}
		}
	}
	return false, 0, 0
}

func runInputParser(envRcv <-chan *environment, envRqst chan<- bool, plyrBuf chan<- inputCommand, stRcv <-chan stateDescriptor, stRqst chan<- stateRequest, stMdfy chan<- stateDescriptor) {
	envRqst <- true
	env := <-envRcv
	
	mouseLeft, mouseRight := false, false
	
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			func() {
				defer func() {recover()}()	//this could maybe print an informative message to a log somewhere
				if event.Ch == 0 && event.Key == termbox.KeyEsc {
					setState(stMdfy, stRcv, stateExit)
				}
				switch state := getState(stRqst, stRcv); state {
				case stateMainMenu:
					if event.Ch == 0 {
						if event.Key == termbox.KeyArrowUp {
							setSubState(stMdfy, stRcv, state, subStateInitialValues[state][int(stateMainMenuSelectorIndex)] + (len(selectorMap[variableContentsKey{state: state, subStateIndex: stateMainMenuSelectorIndex, displayMode: displayMainMenu}]) + getSubState(stRqst, stRcv, state, stateMainMenuSelectorIndex) - 1) % len(selectorMap[variableContentsKey{state: state, subStateIndex: stateMainMenuSelectorIndex, displayMode: displayMainMenu}]), stateMainMenuSelectorIndex)
						}else if event.Key == termbox.KeyArrowDown {
							setSubState(stMdfy, stRcv, state, subStateInitialValues[state][int(stateMainMenuSelectorIndex)] + (getSubState(stRqst, stRcv, state, stateMainMenuSelectorIndex) + 1) % len(selectorMap[variableContentsKey{state: state, subStateIndex: stateMainMenuSelectorIndex, displayMode: displayMainMenu}]), stateMainMenuSelectorIndex)
						}else if event.Key == termbox.KeyEnter {
							switch getSubState(stRqst, stRcv, state, stateMainMenuSelectorIndex) {
							case 0:
								setState(stMdfy, stRcv, stateNewGame)
							case 1:
								//options
							case 2:
								setState(stMdfy, stRcv, stateExit)
							}
						}
					}
				case stateNewGame:
					if event.Ch == 0 {
						if event.Key == termbox.KeyEnter {
							setState(stMdfy, stRcv, stateRunningGame)
							env.paused.Open()
						}
					}
				case stateRunningGame:
					if event.Ch == 0 {
						switch event.Key {
						case termbox.KeyArrowLeft:
							plyrBuf <- initKeyCommand(termbox.KeyArrowLeft)
						case termbox.KeyArrowUp:
							plyrBuf <- initKeyCommand(termbox.KeyArrowUp)
						case termbox.KeyArrowRight:
							plyrBuf <- initKeyCommand(termbox.KeyArrowRight)
						case termbox.KeyArrowDown:
							plyrBuf <- initKeyCommand(termbox.KeyArrowDown)
						case termbox.KeySpace:
							env.paused.Close()
							setState(stMdfy, stRcv, statePausedGame)
						}
					}
				case statePausedGame:
					if event.Ch == 0 {
						if event.Key == termbox.KeyArrowUp {
							setSubState(stMdfy, stRcv, state, subStateInitialValues[state][int(statePausedGameSelectorIndex)] + (len(selectorMap[variableContentsKey{state: state, subStateIndex: statePausedGameSelectorIndex, displayMode: displayPause}]) + getSubState(stRqst, stRcv, state, statePausedGameSelectorIndex) - 1) % len(selectorMap[variableContentsKey{state: state, subStateIndex: statePausedGameSelectorIndex, displayMode: displayPause}]), statePausedGameSelectorIndex)
						}else if event.Key == termbox.KeyArrowDown {
							setSubState(stMdfy, stRcv, state, subStateInitialValues[state][int(statePausedGameSelectorIndex)] + (getSubState(stRqst, stRcv, state, statePausedGameSelectorIndex) + 1) % len(selectorMap[variableContentsKey{state: state, subStateIndex: statePausedGameSelectorIndex, displayMode: displayPause}]), statePausedGameSelectorIndex)
						}else if event.Key == termbox.KeyEnter {
							switch getSubState(stRqst, stRcv, state, statePausedGameSelectorIndex) {
							case 0:
								setState(stMdfy, stRcv, stateRunningGame)
								env.paused.Open()
							case 1:
								setState(stMdfy, stRcv, stateMainMenu)
							}
						}
					}
				case stateGameOver:
					if event.Ch == 0 {
						if event.Key == termbox.KeyEnter {
							switch getSubState(stRqst, stRcv, state, stateGameOverSelectorIndex) {
							case 0:
								setState(stMdfy, stRcv, stateMainMenu)
							}
						}
					}
				}
			}()
		}else if event.Type == termbox.EventMouse {
			if event.Key == termbox.MouseLeft {
				mouseLeft, mouseRight = true, false
			}else if event.Key == termbox.MouseRight {
				mouseLeft, mouseRight = false, true
			}else{
				if event.Key == termbox.MouseRelease {
					if state := getState(stRqst, stRcv); state == stateRunningGame {
						if valid, x, y := mapClickToEnvironment(event, env, state); valid {
							if mouseLeft {
								plyrBuf <- initMouseCommand(termbox.MouseLeft, x, y)
							}else if mouseRight {
								plyrBuf <- initMouseCommand(termbox.MouseRight, x, y)
							}
						}
					}
				}
				mouseLeft, mouseRight = false, false
			}
		}
	}
}