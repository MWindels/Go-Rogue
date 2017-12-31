package main

import "github.com/nsf/termbox-go"

const (
	stateMainMenu uint8 = iota
	stateNewGame
	stateRunningGame
	statePausedGame
	stateExit
)

type game struct {
	state uint8
	
	envSnd chan *environment
	envRqst chan bool
	envMdfy chan bool					//should become it's own type
	
	entSnd chan *entity
	
	stSnd chan uint8
	stRqst chan bool
	stMdfy chan uint8
}

func initGame() game {
	g := game{
		state: stateMainMenu,
		
		envSnd: make(chan *environment),
		envRqst: make(chan bool),
		envMdfy: make(chan bool),			//likewise here
		
		entSnd: make(chan *entity),
		
		stSnd: make(chan uint8),
		stRqst: make(chan bool),
		stMdfy: make(chan uint8),
	}
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
		case <- g.stRqst:
			g.stSnd <- g.state
		case newState := <- g.stMdfy:
			g.state = newState
		}
	}
}