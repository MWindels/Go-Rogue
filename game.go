package main

import "github.com/nsf/termbox-go"

type game struct {
	state uint
	subStates []int
	
	envSnd chan *environment
	envRqst chan bool
	envMdfy chan bool					//should become it's own type
	
	entSnd chan *entity
	
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
		envMdfy: make(chan bool),			//likewise here
		
		entSnd: make(chan *entity),
		
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
	
	g := initGame()
	
	//Testing stuff
	go runEnvController(g.envSnd, g.envRqst, g.envMdfy, g.entSnd)
	go runRenderer(g.envSnd, g.envRqst, g.stSnd, g.stRqst)
	go runInputParser(g.envSnd, g.envRqst, g.envMdfy, g.stSnd, g.stRqst, g.stMdfy)
	go runEntity(g.envSnd, g.envRqst, g.entSnd)
	//Testing stuff
	
	for g.state != stateExit {
		select{
		case req := <- g.stRqst:
			if req.reqType == stateType {
				g.stSnd <- initStateDesc(g.state)
			}else if req.reqType == subStateType {
				if req.subStateIndex < uint(len(g.subStates)) {
					g.stSnd <- initSubStateDesc(g.state, g.subStates[req.subStateIndex], req.subStateIndex)
				}else{
					g.stSnd <- initErrorDesc()
				}
			}else{
				g.stSnd <- initErrorDesc()
			}
		case mod := <- g.stMdfy:
			if mod.state < totalStates {
				if mod.descType == stateType {
					g.state = mod.state
					g.subStates = make([]int, int(totalSubStates[g.state]), int(totalSubStates[g.state]))
					copy(g.subStates, subStateInitialValues[g.state])
				}else if mod.descType == subStateType {
					if mod.state == g.state {
						if mod.subStateIndex < uint(len(g.subStates)) {
							g.subStates[mod.subStateIndex] = mod.subState
						}
					}
				}
			}
		}
	}
}