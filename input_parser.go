package main

import "github.com/nsf/termbox-go"

func runInputParser(envRcv <-chan *environment, envRqst chan<- bool, envMdfy chan<- bool, stRcv <-chan stateDescriptor, stRqst chan<- stateRequest, stMdfy chan<- stateDescriptor) {
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			func(){
				defer func(){recover()}()	//this could maybe print an informative message
				if event.Ch == 0 && event.Key == termbox.KeyEsc {
					stMdfy <- initStateDesc(stateExit)
				}
				switch state := getState(stRqst, stRcv); state {
				case stateMainMenu:
					if event.Ch == 0 {
						if event.Key == termbox.KeyArrowUp {
							stMdfy <- initSubStateDesc(state, subStateInitialValues[state][int(stateMainMenuSelectorIndex)] + (len(selectorMap[variableContentsKey{state: state, subStateIndex: stateMainMenuSelectorIndex, displayMode: displayMainMenu}]) + getSubState(stRqst, stRcv, state, stateMainMenuSelectorIndex) - 1) % len(selectorMap[variableContentsKey{state: state, subStateIndex: stateMainMenuSelectorIndex, displayMode: displayMainMenu}]), stateMainMenuSelectorIndex)
						}else if event.Key == termbox.KeyArrowDown {
							stMdfy <- initSubStateDesc(state, subStateInitialValues[state][int(stateMainMenuSelectorIndex)] + (getSubState(stRqst, stRcv, state, stateMainMenuSelectorIndex) + 1) % len(selectorMap[variableContentsKey{state: state, subStateIndex: stateMainMenuSelectorIndex, displayMode: displayMainMenu}]), stateMainMenuSelectorIndex)
						}else if event.Key == termbox.KeyEnter {
							switch getSubState(stRqst, stRcv, state, stateMainMenuSelectorIndex){
							case 0:
								stMdfy <- initStateDesc(stateNewGame)
							case 1:
								//options
							case 2:
								stMdfy <- initStateDesc(stateExit)
							}
						}
					}
				case stateNewGame:
					if event.Ch == 0 {
						if event.Key == termbox.KeyEnter {
							stMdfy <- initStateDesc(stateRunningGame)
						}
					}
				case stateRunningGame:
					if event.Ch == 0 {
						if event.Key == termbox.KeySpace {
							stMdfy <- initStateDesc(statePausedGame)
						}
					}
				case statePausedGame:
					if event.Ch == 0 {
						if event.Key == termbox.KeyArrowUp {
							stMdfy <- initSubStateDesc(state, subStateInitialValues[state][int(statePausedGameSelectorIndex)] + (len(selectorMap[variableContentsKey{state: state, subStateIndex: statePausedGameSelectorIndex, displayMode: displayPause}]) + getSubState(stRqst, stRcv, state, statePausedGameSelectorIndex) - 1) % len(selectorMap[variableContentsKey{state: state, subStateIndex: statePausedGameSelectorIndex, displayMode: displayPause}]), statePausedGameSelectorIndex)
						}else if event.Key == termbox.KeyArrowDown {
							stMdfy <- initSubStateDesc(state, subStateInitialValues[state][int(statePausedGameSelectorIndex)] + (getSubState(stRqst, stRcv, state, statePausedGameSelectorIndex) + 1) % len(selectorMap[variableContentsKey{state: state, subStateIndex: statePausedGameSelectorIndex, displayMode: displayPause}]), statePausedGameSelectorIndex)
						}else if event.Key == termbox.KeyEnter {
							switch getSubState(stRqst, stRcv, state, statePausedGameSelectorIndex){
							case 0:
								stMdfy <- initStateDesc(stateRunningGame)
							case 1:
								stMdfy <- initStateDesc(stateMainMenu)
							}
						}
					}
				}
			}()
			
			/*if event.Ch == 'n' {
				stMdfy <- initStateDesc((getState(stRqst, stRcv) + 1) % stateExit)
			}else if event.Ch == 'd' {
				func(){
					defer func(){recover()}()		//this could maybe print an informative message
					state := getState(stRqst, stRcv)
					for i := 0; i < int(totalSubStates[state]); i++ {
						subState := getSubState(stRqst, stRcv, state, uint(i))
						stMdfy <- initSubStateDesc(state, (subState + 1) % 5, uint(i))	//this essentially must be done on a case-by-case basis, since a system of substate descriptors can't adequately encapsulate the complete funtionality of substates
					}
				}()
			}*/
		}
	}
}