package main

import "github.com/nsf/termbox-go"

func runInputParser(envRcv <-chan *environment, envRqst chan<- bool, plyrBuf chan<- termbox.Key, stRcv <-chan stateDescriptor, stRqst chan<- stateRequest, stMdfy chan<- stateDescriptor) {
	envRqst <- true
	env := <-envRcv
	
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
							env.paused.Done()
						}
					}
				case stateRunningGame:
					if event.Ch == 0 {
						switch event.Key {
						case termbox.KeyArrowLeft:
							plyrBuf <- termbox.KeyArrowLeft
						case termbox.KeyArrowUp:
							plyrBuf <- termbox.KeyArrowUp
						case termbox.KeyArrowRight:
							plyrBuf <- termbox.KeyArrowRight
						case termbox.KeyArrowDown:
							plyrBuf <- termbox.KeyArrowDown
						case termbox.KeySpace:
							setState(stMdfy, stRcv, statePausedGame)
							env.paused.Add(1)
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
								env.paused.Done()
							case 1:
								setState(stMdfy, stRcv, stateMainMenu)
							}
						}
					}
				}
			}()
		}
	}
}