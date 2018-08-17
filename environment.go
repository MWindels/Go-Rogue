package main

import (
	"sync"
	"math/rand"
	"github.com/nsf/termbox-go"
	"github.com/mwindels/go-rogue/adt"
	"github.com/mwindels/go-rogue/geom"
)

const (
	roomCreateLow float64 = 0.7
	roomCreateHigh float64 = 0.8
	roomDestroyLow float64 = 0.9
	roomDestroyHigh float64 = 1.0
)

type environment struct {
	mutex sync.RWMutex
	width int
	height int
	tiles [][]termbox.Cell
	entities [][](*entity)
	rooms adt.BSPTree
}

func generatePartitions(w, h int, partitionCount int) []geom.Point {
	partitions := make([]geom.Point, partitionCount, partitionCount)
	for i := 0; i < partitionCount; i++ {
		partitions[i] = geom.InitPoint(rand.Float64() * float64(w), rand.Float64() * float64(h))
	}
	return partitions
}

func initEnvironment(w, h int) environment {
	env := environment{
		mutex: sync.RWMutex{},
		width: w,
		height: h,
		tiles: make([][]termbox.Cell, w),
		entities: make([][](*entity), w),
		rooms: adt.InitBSPTree(geom.InitRectangle(0.0, 0.0, float64(w), float64(h)), generatePartitions(w, h, w * h)),		//w * h is just a stand-in for a more complicated partition count function
	}
	for row := 0; row < int(w); row++ {
		env.tiles[row] = make([]termbox.Cell, h)
		env.entities[row] = make([](*entity), h)
	}
	return env
}

func createRoomPredicate(relativeDepth float64) bool {
	return roomCreateLow <= relativeDepth && relativeDepth <= roomCreateHigh && rand.Float64() * (roomCreateHigh - roomCreateLow) <= 0.5 * (relativeDepth - roomCreateLow)
}

func destroyRoomPredicate(relativeDepth float64) bool {
	return roomDestroyLow <= relativeDepth && relativeDepth <= roomDestroyHigh && rand.Float64() * (roomDestroyHigh - roomDestroyLow) <= 0.75 * (relativeDepth - roomDestroyLow)
}

func (e *environment) generateRooms(randomizations int) {
	for i := 0; i < randomizations; i++ {
		e.rooms.RandomizeTraversability(createRoomPredicate, destroyRoomPredicate)
	}
	
	var stk adt.Stack
	stk.Push(e.rooms.Root())
	for !(stk.IsEmpty()) {
		node, valid := stk.Pop().(*adt.BSPNode)
		if !valid {
			panic("Popped an element from the tile-filling stack that wasn't a BSPNode pointer!")
		}
		if node.Traversability() != adt.SemiOpen {
			for x := int(node.Area().UpperLeft().X); x < int(node.Area().UpperRight().X); x++ {
				for y := int(node.Area().UpperLeft().Y); y < int(node.Area().LowerLeft().Y); y++ {
					if node.Traversability() == adt.Open {
						e.tiles[x][y] = termbox.Cell{Ch: '.', Fg: termbox.ColorDefault, Bg: termbox.ColorBlack}
					}else if node.Traversability() == adt.Closed {
						e.tiles[x][y] = termbox.Cell{Ch: '#', Fg: termbox.ColorDefault, Bg: termbox.ColorBlack}
					}
				}
			}
		}else{	//A semi-open node implies that the node is not a leaf.
			stk.Push(node.Left())
			stk.Push(node.Right())
		}
	}
}

func generateEnvironment(e environment) environment {
	e.generateRooms(1)
	return e
}

func runEnvController(envSnd chan<- *environment, envRqst <-chan bool, envMdfy <-chan bool, entRcv <-chan *entity) {			//remember, change the type of envMdfy some time
	env := generateEnvironment(initEnvironment(100, 100))
	
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