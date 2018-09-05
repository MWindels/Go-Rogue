package main

import (
	"time"
	"sync"
	"math"
	"math/rand"
	"github.com/nsf/termbox-go"
	"github.com/mwindels/go-rogue/adt"
	"github.com/mwindels/go-rogue/geom"
)

type entity struct {
	deathNotify chan bool
	targeted sync.WaitGroup
	
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
		deathNotify: make(chan bool, 1),
		targeted: sync.WaitGroup{},
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
			}else{
				env.entities[ent.x + xOffset][ent.y + yOffset].kill()	//could then remove the entity, but then it may appear as if entities carry out attacks after they die.
			}
		}
	}
	env.mutex.Unlock()
	return moved
}

func (ent entity) findPath(env *environment, x, y int, useLock bool) []geom.Point {
	if useLock {
		env.mutex.RLock()
		defer env.mutex.RUnlock()
	}
	
	var pathQueue adt.PriorityQueue
	previousTile := make(map[geom.Point]geom.Point)
	start, end := geom.InitPoint(float64(ent.x), float64(ent.y)), geom.InitPoint(float64(x), float64(y))
	
	pathQueue.Insert([2]geom.Point{start, start}, geom.PointDistance(start, end))
	for !(pathQueue.IsEmpty()) {
		temp, tilePriority := pathQueue.Extract()
		tilePair, valid := temp.([2]geom.Point)
		if !valid {
			panic("Extracted an element from the path-finding priority queue which wasn't an array of two Points!")
		}
		
		if _, exists := previousTile[tilePair[0]]; !exists {
			previousTile[tilePair[0]] = tilePair[1]
			if geom.PointsEqual(tilePair[0], end) {
				pathQueue.Clear()
			}else{
				for i := 0; i < 4; i++ {
					nextTile := geom.InitPoint(tilePair[0].X + neighbourPoints[2 * i + 1].X, tilePair[0].Y + neighbourPoints[2 * i + 1].Y)
					if geom.RectangleContainsLowerInclusive(env.rooms.Root().Area(), nextTile) {
						if env.tiles[int(nextTile.X)][int(nextTile.Y)] == floorTile /*&& env.entities[int(nextTile.X)][int(nextTile.Y)] == nil*/ {
							pathQueue.Insert([2]geom.Point{nextTile, tilePair[0]}, tilePriority - geom.PointDistance(tilePair[0], end) + geom.PointDistance(tilePair[0], nextTile) + geom.PointDistance(nextTile, end))
						}
					}
				}
			}
		}
	}
	
	if _, exists := previousTile[end]; exists {
		if !(geom.PointsEqual(start, end)) {
			backwardsPath := []geom.Point{end}
			for tile := previousTile[end]; !(geom.PointsEqual(tile, start)); tile = previousTile[tile] {
				backwardsPath = append(backwardsPath, tile)
			}
			forwardsPath := []geom.Point{start}
			for i := len(backwardsPath) - 1; i >= 0; i-- {
				forwardsPath = append(forwardsPath, backwardsPath[i])
			}
			return forwardsPath
		}else{
			return []geom.Point{}
		}
	}else{
		return []geom.Point{}
	}
}

func (ent *entity) findTarget(env *environment, x, y int, maxDistance float64) *entity {
	env.mutex.RLock()
	defer env.mutex.RUnlock()
	
	var nearestQueue adt.PriorityQueue
	checkedTiles := make(map[geom.Point]bool)
	start := geom.InitPoint(float64(x), float64(y))
	
	nearestQueue.Insert(start, 0.0)
	for !(nearestQueue.IsEmpty()) {
		temp, tileDistance := nearestQueue.Extract()
		tile, valid := temp.(geom.Point)
		if !valid {
			panic("Extracted an element from the entity-finding priority queue that wasn't a Point!")
		}
		
		if tileDistance <= maxDistance {
			if _, exists := checkedTiles[tile]; !exists {
				checkedTiles[tile] = true
				if env.entities[int(tile.X)][int(tile.Y)] != nil && env.entities[int(tile.X)][int(tile.Y)] != ent {
					env.entities[int(tile.X)][int(tile.Y)].targeted.Add(1)
					return env.entities[int(tile.X)][int(tile.Y)]
				}else{
					for i := 0; i < 4; i++ {
						nextTile := geom.InitPoint(tile.X + neighbourPoints[2 * i + 1].X, tile.Y + neighbourPoints[2 * i + 1].Y)
						if geom.RectangleContainsLowerInclusive(env.rooms.Root().Area(), nextTile) {
							nearestQueue.Insert(nextTile, tileDistance + geom.PointDistance(tile, nextTile))
						}
					}
				}
			}
		}
	}
	
	return nil
}

func (ent entity) findPathToTarget(env *environment, target *entity) ([]geom.Point, bool) {
	env.mutex.RLock()
	defer env.mutex.RUnlock()
	
	if env.entities[target.x][target.y] == target {	//Make sure the target still exists in the environment.
		return ent.findPath(env, target.x, target.y, false)[:int(math.Max(geom.PointDistance(geom.InitPoint(float64(ent.x), float64(ent.y)), geom.InitPoint(float64(target.x), float64(target.y))) / (float64(ent.wait) / float64(target.wait) + 1), 2))], true
	}else{
		return []geom.Point{}, false
	}
}

func (ent *entity) wander(env *environment) {
	if rand.Intn(2) == 0 {
		if rand.Intn(2) == 0 {
			ent.move(env, 1, 0)
		}else{
			ent.move(env, -1, 0)
		}
	}else{
		if rand.Intn(2) == 0 {
			ent.move(env, 0, 1)
		}else{
			ent.move(env, 0, -1)
		}
	}
}

func (ent *entity) kill() {
	defer func() {recover()}()	//In case deathNotify has already been closed.
	select{
	case ent.deathNotify <- true:	//This will only block if something is in the channel, and it will panic if the channel has been closed.
		close(ent.deathNotify)
		return
	default:
		return
	}
}

func (ent *entity) die(env *environment) {
	env.mutex.Lock()
	env.entities[ent.x][ent.y] = nil	//maybe drop a corpse item afterwards?
	env.mutex.Unlock()
}

func runEntity(envRcv <-chan *environment, envRqst chan<- bool, ent entity) {
	envRqst <- true
	env := <-envRcv
	
	if !ent.spawn(env) {	//if this succeeds, in order to access the mutable fields in ent you must lock env's mutex(es)
		return
	}
	
	alive := true
	cooldown := time.NewTicker(ent.wait)
	for alive {
		env.paused.Wait()
		select{
		case <-ent.deathNotify:
			ent.die(env)
			alive = false
		case <-cooldown.C:
			env.paused.Wait()
			ent.wander(env)		//perform in separate thread?
		}
	}
	cooldown.Stop()
	
	ent.targeted.Wait()
}

func runPlayer(envRcv <-chan *environment, envRqst chan<- bool, plyrBuf <-chan inputCommand, stRcv <-chan stateDescriptor, stMdfy chan<- stateDescriptor, plyr entity) {
	envRqst <- true
	env := <-envRcv
	
	if plyr.findSpawn(env) {
		env.player = &plyr
	}else{
		for range plyrBuf {}
		return
	}
	
	alive := true
	pathPosition := 1
	var currentPath []geom.Point
	var currentTarget *entity = nil
	cooldown := time.NewTicker(plyr.wait)
	
	for alive {
		env.paused.Wait()
		select{
		case <-plyr.deathNotify:
			plyr.die(env)
			alive = false
		case cmd := <-plyrBuf:
			if cmd.commandType == commandKey {
				switch cmd.key {
				case termbox.KeyArrowLeft:
					plyr.move(env, -1, 0)
				case termbox.KeyArrowUp:
					plyr.move(env, 0, -1)
				case termbox.KeyArrowRight:
					plyr.move(env, 1, 0)
				case termbox.KeyArrowDown:
					plyr.move(env, 0, 1)
				}
			}else if cmd.commandType == commandMouse {
				if geom.RectangleContainsLowerInclusive(env.rooms.Root().Area(), geom.InitPoint(float64(cmd.x), float64(cmd.y))) {
					switch cmd.key {
					case termbox.MouseLeft:
						pathPosition = 1
						currentPath = plyr.findPath(env, cmd.x, cmd.y, true)
						if currentTarget != nil {
							currentTarget.targeted.Done()
							currentTarget = nil
						}
					case termbox.MouseRight:
						if currentTarget != nil {
							currentTarget.targeted.Done()
						}
						currentTarget = plyr.findTarget(env, cmd.x, cmd.y, 5.0)
						if currentTarget != nil {
							if temp, valid := plyr.findPathToTarget(env, currentTarget); valid {	//construct initial path to target
								pathPosition = 1
								currentPath = temp
							}else{
								currentTarget.targeted.Done()
								currentTarget = nil
							}
						}
					}
				}
			}
		case <-cooldown.C:
			if pathPosition < len(currentPath) {	//add check to make sure you are where you think you are (allowing for at least one tile error due to unlock before moving)?
				if plyr.move(env, int(currentPath[pathPosition].X) - int(currentPath[pathPosition - 1].X), int(currentPath[pathPosition].Y) - int(currentPath[pathPosition - 1].Y)) {
					pathPosition += 1
				}else{
					pathPosition = len(currentPath)
				}
			}
			if pathPosition >= len(currentPath) && currentTarget != nil {
				if temp, valid := plyr.findPathToTarget(env, currentTarget); valid {	//re-path after initial (or subsequent) path to target runs out
					pathPosition = 1
					currentPath = temp
				}else{
					currentTarget.targeted.Done()
					currentTarget = nil
				}
			}
		}
	}
	
	if currentTarget != nil {
		currentTarget.targeted.Done()
	}
	
	env.paused.Close()
	setState(stMdfy, stRcv, stateGameOver)	//This won't panic, because it's sending a valid state.
	for range plyrBuf {}
	plyr.targeted.Wait()
}