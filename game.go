package main

import (
	"time"
	"math/rand"
	"github.com/nsf/termbox-go"
)

type game struct {
	state uint
	subStates []int
	
	envSnd chan *environment
	envRqst chan bool
	
	plyrBuf chan inputCommand
	
	stSnd chan stateDescriptor
	stRqst chan stateRequest
	stMdfy chan stateDescriptor
}

func initGame() game {
	g := game{
		state: stateMainMenu,
		subStates: make([]int, int(totalSubStates[stateMainMenu]), int(totalSubStates[stateMainMenu])),
		
		envSnd: make(chan *environment),
		envRqst: make(chan bool),
		
		plyrBuf: make(chan inputCommand, 1028),
		
		stSnd: make(chan stateDescriptor),
		stRqst: make(chan stateRequest),
		stMdfy: make(chan stateDescriptor),
	}
	copy(g.subStates, subStateInitialValues[stateMainMenu])
	return g
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	
	rand.Seed(time.Now().UTC().UnixNano())
	
	g := initGame()
	
	go runEnvController(g.envSnd, g.envRqst)
	go runRenderer(g.envSnd, g.envRqst, g.stSnd, g.stRqst)
	go runInputParser(g.envSnd, g.envRqst, g.plyrBuf, g.stSnd, g.stRqst, g.stMdfy)
	go runPlayer(g.envSnd, g.envRqst, g.plyrBuf, g.stSnd, g.stMdfy, initEntity('@', termbox.ColorDefault, 100, 0, 0))
	
	for g.state != stateExit {
		select{
		case req := <-g.stRqst:
			if req.reqType == stateType {
				g.stSnd <- initStateDesc(g.state)
				break
			}else if req.reqType == subStateType {
				if req.subStateIndex < uint(len(g.subStates)) {
					g.stSnd <- initSubStateDesc(g.state, g.subStates[req.subStateIndex], req.subStateIndex)
					break
				}
			}
			g.stSnd <- initErrorDesc()
		case mod := <-g.stMdfy:
			if mod.state < totalStates {
				if mod.descType == stateType {
					g.state = mod.state
					g.subStates = make([]int, int(totalSubStates[g.state]), int(totalSubStates[g.state]))
					copy(g.subStates, subStateInitialValues[g.state])
					g.stSnd <- mod
					break
				}else if mod.descType == subStateType {
					if mod.state == g.state {
						if mod.subStateIndex < uint(len(g.subStates)) {
							g.subStates[mod.subStateIndex] = mod.subState
							g.stSnd <- mod
							break
						}
					}
				}
			}
			g.stSnd <- initErrorDesc()
		}
	}
}