package main

import (
	"time"
	"math/rand"
	"github.com/nsf/termbox-go"
	"github.com/mwindels/go-rogue/adt"
	"github.com/mwindels/go-rogue/geom"
)

type entity struct {
	//Immutable
	symbol rune
	color termbox.Attribute
	wait time.Duration
	
	//Mutable
	x int
	y int
}

func initEntity(s rune, col termbox.Attribute, ms uint, x, y int) entity {
	ent := entity{
		symbol: s,
		color: col,
		wait: time.Duration(ms) * time.Millisecond,
		x: x,
		y: y,
	}
	return ent
}

func (ent *entity) spawn(env *environment) bool {
	spawned := false
	if geom.RectangleContainsLowerInclusive(env.rooms.Root().Area(), geom.InitPoint(float64(ent.x), float64(ent.y))) {
		if env.tiles[ent.x][ent.y] == floorTile {
			env.mutex.Lock()	//Only need to lock here because ent is not yet a part of env.
			if env.entities[ent.x][ent.y] == nil {
				env.entities[ent.x][ent.y] = ent
				spawned = true
			}
			env.mutex.Unlock()
		}
	}
	return spawned
}

func (ent *entity) findSpawn(env *environment) bool {
	var stk adt.Stack
	func() {
		defer func() {recover()}()
		stk.Push(env.rooms.Root())
	}()
	for !(stk.IsEmpty()) {
		node, valid := stk.Pop().(adt.BSPNode)
		if !valid {
			panic("Popped an element from the spawn-finding stack that wasn't a BSPNode!")
		}
		if node.Traversability() == adt.SemiOpen {
			stk.Push(node.Right())
			stk.Push(node.Left())
		}else if node.Traversability() == adt.Open {
			ent.x = int(node.Area().UpperLeft().X + rand.Float64() * node.Area().Width())
			ent.y = int(node.Area().UpperLeft().Y + rand.Float64() * node.Area().Height())
			if !ent.spawn(env) {
				for ent.x = int(node.Area().UpperLeft().X); ent.x < int(node.Area().UpperRight().X); ent.x++ {
					for ent.y = int(node.Area().UpperLeft().Y); ent.y < int(node.Area().LowerLeft().Y); ent.y++ {
						if ent.spawn(env) {
							stk.Clear()
							return true
						}
					}
				}
			}else{
				stk.Clear()
				return true
			}
		}
	}
	return false
}

func (ent *entity) move(env *environment, xOffset, yOffset int) bool {
	moved := false
	env.mutex.Lock()
	if geom.RectangleContainsLowerInclusive(env.rooms.Root().Area(), geom.InitPoint(float64(ent.x + xOffset), float64(ent.y + yOffset))) {
		if env.tiles[ent.x + xOffset][ent.y + yOffset] == floorTile {
			if env.entities[ent.x + xOffset][ent.y + yOffset] == nil {
				env.entities[ent.x][ent.y] = nil
				ent.x += xOffset
				ent.y += yOffset
				env.entities[ent.x][ent.y] = ent
				moved = true
			}
		}
	}
	env.mutex.Unlock()
	return moved
}

func (ent *entity) wander(env *environment) {
	env.mutex.RLock()
	x := ent.x
	y := ent.y
	env.mutex.RUnlock()
	
	offset := neighbourPoints[2 * rand.Intn(4) + 1]
	nextTile := geom.InitPoint(float64(x) + offset.X, float64(y) + offset.Y)
	if geom.RectangleContainsLowerInclusive(env.rooms.Root().Area(), nextTile) {
		if env.tiles[int(nextTile.X)][int(nextTile.Y)] == floorTile {
			env.mutex.Lock()
			//if not dead!
			if ent.x == x && ent.y == y {	//optimistic lock
				if env.entities[int(nextTile.X)][int(nextTile.Y)] == nil {
					env.entities[ent.x][ent.y] = nil
					ent.x = int(nextTile.X)
					ent.y = int(nextTile.Y)
					env.entities[ent.x][ent.y] = ent
				}
			}
			env.mutex.Unlock()
		}
	}
}

func runEntity(envRcv <-chan *environment, envRqst chan<- bool, ent entity) {
	envRqst <- true
	env := <-envRcv
	
	if !ent.spawn(env) {	//if this succeeds, in order to access the mutable fields in ent you must lock env's mutex(es)
		return
	}
	
	for /*not dead*/ {	//dying removes ent from env
		env.paused.Wait()
		time.Sleep(ent.wait)
		env.paused.Wait()
		
		ent.wander(env)
	}
}

func runPlayer(envRcv <-chan *environment, envRqst chan<- bool, plyrBuf <-chan termbox.Key, plyr entity) {
	envRqst <- true
	env := <-envRcv
	
	if plyr.findSpawn(env) {
		env.player = &plyr
	}else{
		for range plyrBuf {}
		return
	}
	
	for /*not dead*/ {
		env.paused.Wait()
		select{
		case cmd := <-plyrBuf:
			switch cmd {
			case termbox.KeyArrowLeft:
				plyr.move(env, -1, 0)
			case termbox.KeyArrowUp:
				plyr.move(env, 0, -1)
			case termbox.KeyArrowRight:
				plyr.move(env, 1, 0)
			case termbox.KeyArrowDown:
				plyr.move(env, 0, 1)
			}
		}
	}
}