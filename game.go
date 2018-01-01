package main

import "github.com/nsf/termbox-go"

const (
	stateMainMenu uint8 = iota
	stateNewGame
	stateRunningGame
	statePausedGame
	stateExit
	totalStates
)

var (
	subStatesByState [totalStates]uint = [totalStates]uint{1, 1, 0, 1, 0}
)

type subStateModification struct {
	index uint
	value int
}

type game struct {
	state uint8
	subState []int
	
	envSnd chan *environment
	envRqst chan bool
	envMdfy chan bool					//should become it's own type
	
	entSnd chan *entity
	
	stSnd chan uint8
	stRqst chan bool
	stMdfy chan uint8
	
	sstSnd chan int
	sstRqst chan uint
	sstMdfy chan subStateModification
}

func initSubStateModification(ind uint, val int) subStateModification {
	ssm := subStateModification{
		index: ind,
		value: val,
	}
	return ssm
}

func initGame() game {
	g := game{
		state: stateMainMenu,
		subState: make([]int, subStatesByState[stateMainMenu], subStatesByState[stateMainMenu]),
		
		envSnd: make(chan *environment),
		envRqst: make(chan bool),
		envMdfy: make(chan bool),			//likewise here
		
		entSnd: make(chan *entity),
		
		stSnd: make(chan uint8),
		stRqst: make(chan bool),
		stMdfy: make(chan uint8),
		
		sstSnd: make(chan int),
		sstRqst: make(chan uint),
		sstMdfy: make(chan subStateModification),
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
	go runRenderer(g.envSnd, g.envRqst, g.stSnd, g.stRqst, g.sstSnd, g.sstRqst)
	go runInputParser(g.envSnd, g.envRqst, g.envMdfy, g.stSnd, g.stRqst, g.stMdfy, g.sstSnd, g.sstRqst, g.sstMdfy)
	go runEntity(g.envSnd, g.envRqst, g.entSnd)
	//Testing stuff
	
	for g.state != stateExit {
		select{
		case <- g.stRqst:
			g.stSnd <- g.state
		case newState := <- g.stMdfy:
			if 0 <= newState && newState < totalStates {
				g.state = newState
				g.subState = make([]int, subStatesByState[g.state], subStatesByState[g.state])
			}
		case index := <- g.sstRqst:
			if index < subStatesByState[g.state] {
				g.sstSnd <- g.subState[index]
			}
		case modification := <- g.sstMdfy:
			if modification.index < subStatesByState[g.state] {
				g.subState[modification.index] = modification.value
			}
		}
	}
}