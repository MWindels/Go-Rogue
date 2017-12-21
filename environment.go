package main

import "sync"

type environment struct {
	mutex sync.RWMutex
	width uint32
	height uint32
	tiles [][]rune
	entities [][](*entity)
}

func initEnvironment(w uint32, h uint32) environment {
	env := environment{
		mutex: sync.RWMutex{},
		width: w,
		height: h,
		tiles: make([][]rune, w),
		entities: make([][](*entity), w),
	}
	for row := 0; row < int(w); row++ {
		env.tiles[row] = make([]rune, h)
		env.entities[row] = make([](*entity), h)
	}
	return env
}

func runEnvController(envSnd chan<- *environment, envRqst <-chan bool, envMdfy <-chan bool, entRcv <-chan *entity) {			//remember, change the type of envMdfy some time
	//generate/load here first
	env := initEnvironment(12, 8)
	
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