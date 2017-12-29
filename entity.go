package main

import "github.com/nsf/termbox-go"

type entity struct {
	symbol rune
	color termbox.Attribute
	x int
	y int
}

func initEntity(s rune, col termbox.Attribute, x int, y int) entity {
	ent := entity{
		symbol: s,
		color: col,
		x: x,
		y: y,
	}
	return ent
}

func runEntity(envRcv <-chan *environment, envRqst chan<- bool, entSnd chan<- *entity) {
	ent := initEntity('@', 0, 1, 2)
	
	entSnd <- &ent
	envRqst <- true
	env := <-envRcv
	
	//Testing
	env.mutex.RLock()
	env.mutex.RUnlock()
	//Testing
	
	for {}
}