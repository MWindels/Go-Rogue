package main

import (
	"sync"
	"math/rand"
	"github.com/nsf/termbox-go"
)

type environment struct {
	mutex sync.RWMutex
	width int
	height int
	tiles [][]termbox.Cell
	entities [][](*entity)
}

func initEnvironment(w int, h int) environment {
	env := environment{
		mutex: sync.RWMutex{},
		width: w,
		height: h,
		tiles: make([][]termbox.Cell, w),
		entities: make([][](*entity), w),
	}
	for row := 0; row < int(w); row++ {
		env.tiles[row] = make([]termbox.Cell, h)
		env.entities[row] = make([](*entity), h)
	}
	return env
}

//Totally random.  Just temporary for now.
func generateEnvironment(e environment) environment {
	for x := 0; x < e.width; x++ {
		for y := 0; y < e.height; y++ {
			e.tiles[x][y] = termbox.Cell{Ch: rune(rand.Int() % 10 + 48), Fg: termbox.Attribute(rand.Int() % 7 + 2), Bg: termbox.Attribute(rand.Int() % 7 + 2)}
		}
	}
	return e
}

func runEnvController(envSnd chan<- *environment, envRqst <-chan bool, envMdfy <-chan bool, entRcv <-chan *entity) {			//remember, change the type of envMdfy some time
	env := generateEnvironment(initEnvironment(50, 50))
	
	for {
		select{
			case <-envRqst:
				envSnd <- &env
			case modification := <-envMdfy:
				env.mutex.Lock()
				modification = !modification //Testing
				env.mutex.Unlock()
			case newEnt := <-entRcv:
				env.mutex.Lock()
				env.entities[newEnt.x][newEnt.y] = newEnt			//can't have two entities in the same place
				env.mutex.Unlock()
		}
	}
}