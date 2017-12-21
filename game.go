package main

import "github.com/nsf/termbox-go"

const (
	mainMenu uint8 = iota
)

type game struct {
	gameState uint8
	envSnd chan *environment
	envRqst chan bool
	envMdfy chan bool					//should become it's own type
	entSnd chan *entity
}

func initGame() game {
	g := game{
		gameState: mainMenu,
		envSnd: make(chan *environment),
		envRqst: make(chan bool),
		envMdfy: make(chan bool),			//likewise here
		entSnd: make(chan *entity),
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
	
	go runEnvController(g.envSnd, g.envRqst, g.envMdfy, g.entSnd)
	go runRenderer(g.envSnd, g.envRqst)
	go runEntity(g.envSnd, g.envRqst, g.entSnd)
	
	ch := 'f'
	for ch != 'e'{
		ch = termbox.PollEvent().Ch
	}
}