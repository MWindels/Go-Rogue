package main

import (
	"time"
	"github.com/nsf/termbox-go"
)

func drawScreen(env *environment) {
	env.mutex.RLock()
	for x := 0; x < int(env.width); x++ {
		for y := 0; y < int(env.height); y++ {
			if env.entities[x][y] != nil{
				termbox.SetCell(x, y, env.entities[x][y].symbol, 0, 0)
			}else{
				termbox.SetCell(x, y, env.tiles[x][y], 0, 0)
			}
		}
	}
	env.mutex.RUnlock()
	
	err := termbox.Flush()
	if err != nil {
		panic(err)
	}
}

func runRenderer(envRcv <-chan *environment, envRqst chan<- bool) {
	envRqst <- true
	env := <- envRcv
	
	for {
		time.Sleep(time.Second)
		drawScreen(env)
	}
}