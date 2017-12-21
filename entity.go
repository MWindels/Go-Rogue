package main

type entity struct {
	symbol rune
	x uint32
	y uint32
}

func initEntity(s rune, x uint32, y uint32) entity {
	ent := entity{
		symbol: s,
		x: x,
		y: y,
	}
	return ent
}

func runEntity(envRcv <-chan *environment, envRqst chan<- bool, entSnd chan<- *entity) {
	ent := initEntity('@', 1, 2)
	
	entSnd <- &ent
	envRqst <- true
	env := <-envRcv
	
	//Testing
	env.mutex.RLock()
	env.mutex.RUnlock()
	//Testing
	
	for {}
}