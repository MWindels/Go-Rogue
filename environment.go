package main

import (
	"sync"
	"math"
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

var (
	floorTile termbox.Cell = termbox.Cell{Ch: ' ', Fg: termbox.ColorDefault, Bg: termbox.ColorBlack}
	wallTile termbox.Cell = termbox.Cell{Ch: '░', Fg: termbox.ColorDefault, Bg: termbox.ColorBlack}
)

var (
	neighbourPoints [8]geom.Point = [8]geom.Point{
		geom.InitPoint(-1.0, -1.0),
		geom.InitPoint(0.0, -1.0),
		geom.InitPoint(1.0, -1.0),
		geom.InitPoint(1.0, 0.0),
		geom.InitPoint(1.0, 1.0),
		geom.InitPoint(0.0, 1.0),
		geom.InitPoint(-1.0, 1.0),
		geom.InitPoint(-1.0, 0.0),
	}
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
		partitions[i] = geom.InitPoint(math.Trunc(rand.Float64() * float64(w)), math.Trunc(rand.Float64() * float64(h)))
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
	func() {
		defer func() {recover()}()
		stk.Push(e.rooms.Root())
	}()
	for !(stk.IsEmpty()) {
		node, valid := stk.Pop().(adt.BSPNode)
		if !valid {
			panic("Popped an element from the tile-filling stack that wasn't a BSPNode!")
		}
		if node.Traversability() != adt.SemiOpen {
			for x := int(node.Area().UpperLeft().X); x < int(node.Area().UpperRight().X); x++ {
				for y := int(node.Area().UpperLeft().Y); y < int(node.Area().LowerLeft().Y); y++ {
					if node.Traversability() == adt.Open {
						e.tiles[x][y] = floorTile
					}else if node.Traversability() == adt.Closed {
						e.tiles[x][y] = wallTile
					}
				}
			}
		}else{	//A semi-open node implies that the node is not a leaf, so no need to recover from Left() or Right().
			stk.Push(node.Left())
			stk.Push(node.Right())
		}
	}
}

func nearestRoomPairs(leftRooms, rightRooms []geom.Rectangle) ([][2]geom.Rectangle, float64) {
	distance := math.Inf(1)
	var nearest [][2]geom.Rectangle
	for _, current := range leftRooms {
		currentNearest, currentDistance := geom.NearestRectangles(current, rightRooms...)
		if float32(currentDistance) <= float32(distance) {		//cuts them down to fewer bits to reduce the effects of round-off error
			if float32(currentDistance) < float32(distance) {
				distance = currentDistance
				nearest = make([][2]geom.Rectangle, 0, len(currentNearest))
			}
			for i := 0; i < len(currentNearest); i++ {
				nearest = append(nearest, [2]geom.Rectangle{current, currentNearest[i]})
			}
		}
	}
	return nearest, distance
}

func twoSegmentHall(left, right geom.Rectangle, partitionDimension uint) []geom.Point {
	var path []geom.Point
	var start, end geom.Point
	if partitionDimension % 2 == 0 {	//partition is in the x dimension.
		start = geom.InitPoint(left.UpperRight().X - math.Trunc(rand.Float64() * (adt.MinPartitionWidth - 2.0)) - 2.0, left.UpperLeft().Y + math.Trunc(rand.Float64() * (left.Height() - 2.0)) + 1.0)
		end = geom.InitPoint(right.UpperLeft().X + math.Trunc(rand.Float64() * (adt.MinPartitionWidth - 2.0)) + 1.0, right.UpperLeft().Y + math.Trunc(rand.Float64() * (right.Height() - 2.0)) + 1.0)
	}else{	//partition is in the y dimension.
		start = geom.InitPoint(left.UpperLeft().X + math.Trunc(rand.Float64() * (left.Width() - 2.0)) + 1.0, left.LowerLeft().Y - math.Trunc(rand.Float64() * (adt.MinPartitionHeight - 2.0)) - 2.0)
		end = geom.InitPoint(right.UpperLeft().X + math.Trunc(rand.Float64() * (right.Width() - 2.0)) + 1.0, right.UpperLeft().Y + math.Trunc(rand.Float64() * (adt.MinPartitionHeight - 2.0)) + 1.0)
	}
	//start := geom.InitPoint(left.UpperLeft().X + math.Trunc(rand.Float64() * (left.Width() - 2.0)) + 1.0, left.UpperLeft().Y + math.Trunc(rand.Float64() * (left.Height() - 2.0)) + 1.0)
	//end := geom.InitPoint(right.UpperLeft().X + math.Trunc(rand.Float64() * (right.Width() - 2.0)) + 1.0, right.UpperLeft().Y + math.Trunc(rand.Float64() * (right.Height() - 2.0)) + 1.0)
	current := geom.InitPoint(start.X, start.Y)
	moveX := func () {
		for current.X != end.X {
			path = append(path, current)
			if start.X < end.X {
				current.X += 1
			}else{
				current.X -= 1
			}
		}
	}
	moveY := func() {
		for current.Y != end.Y {
			path = append(path, current)
			if start.Y < end.Y {
				current.Y += 1
			}else{
				current.Y -= 1
			}
		}
	}
	if rand.Intn(2) == 0 {
		moveX()
		moveY()
	}else{
		moveY()
		moveX()
	}
	return path
}

func threeSegmentHall(left, right geom.Rectangle, partitionDimension uint) []geom.Point {
	var path []geom.Point
	start := geom.InitPoint(left.UpperLeft().X + math.Trunc(rand.Float64() * (left.Width() - 2.0)) + 1.0, left.UpperLeft().Y + math.Trunc(rand.Float64() * (left.Height() - 2.0)) + 1.0)
	end := geom.InitPoint(right.UpperLeft().X + math.Trunc(rand.Float64() * (right.Width() - 2.0)) + 1.0, right.UpperLeft().Y + math.Trunc(rand.Float64() * (right.Height() - 2.0)) + 1.0)
	current := geom.InitPoint(start.X, start.Y)
	if partitionDimension % 2 == 0 {	//partition is in the x dimension, just like in the BSPTree.
		cornerPoint := geom.InitPoint(left.UpperRight().X + math.Trunc(rand.Float64() * (right.UpperLeft().X - left.UpperRight().X - 2.0)) + 1.0, start.Y)
		for !geom.PointsEqual(current, cornerPoint) {
			path = append(path, current, cornerPoint)
			current.X += 1
		}
		for current.Y != end.Y {
			path = append(path, current)
			if start.Y < end.Y {
				current.Y += 1
			}else{
				current.Y -= 1
			}
		}
		for !geom.PointsEqual(current, end) {
			path = append(path, current)
			current.X += 1
		}
	}else{	//partition is in the y dimension (like BSPTree).
		cornerPoint := geom.InitPoint(start.X, left.LowerLeft().Y + math.Trunc(rand.Float64() * (right.UpperLeft().Y - left.LowerLeft().Y - 2.0)) + 1.0)
		for !geom.PointsEqual(current, cornerPoint) {
			path = append(path, current)
			current.Y += 1
		}
		for current.X != end.X {
			path = append(path, current)
			if start.X < end.X {
				current.X += 1
			}else{
				current.X -= 1
			}
		}
		for !geom.PointsEqual(current, end) {
			path = append(path, current)
			current.Y += 1
		}
	}
	return path
}

func (e *environment) connectRooms(node adt.BSPNode, depth uint) []geom.Rectangle {
	if node.Traversability() == adt.Closed {
		return []geom.Rectangle{}
	}else if node.Traversability() == adt.Open {
		return []geom.Rectangle{node.Area()}
	}
	
	leftRooms := e.connectRooms(node.Left(), depth + 1)		//no need to recover, since the node is implicitly semi-open, hence not a leaf.
	rightRooms := e.connectRooms(node.Right(), depth + 1)	//likewise.
	if len(leftRooms) > 0 && len(rightRooms) > 0 {
		var hall []geom.Point
		nearestLeft, leftDistance := geom.NearestRectangles(node.Right().Area(), leftRooms...)
		nearestRight, rightDistance := geom.NearestRectangles(node.Left().Area(), rightRooms...)
		if leftDistance == 0.0 && rightDistance == 0.0 {
			nearestPairs, nearestDistance := nearestRoomPairs(nearestLeft, nearestRight)
			randomPair := nearestPairs[rand.Intn(len(nearestPairs))]
			if nearestDistance > 0.0 {	//may need to prefer room pairs which share edges over those which only share corners (but both have a distance of 0).
				hall = twoSegmentHall(randomPair[0], randomPair[1], depth % 2)
			}
		}else{
			hall = threeSegmentHall(nearestLeft[rand.Intn(len(nearestLeft))], nearestRight[rand.Intn(len(nearestRight))], depth % 2)
		}
		
		start, end := 0, len(hall)
		for i, tile := range hall {
			/*for j := 0; j < 4; j++ {
				neighbour := geom.InitPoint(tile.X + neighbourPoints[2 * j + 1].X, tile.Y + neighbourPoints[2 * j + 1].Y)
				if geom.RectangleContainsSemiInclusive(node.Area(), neighbour) {
					if e.tiles[int(neighbour.X)][int(neighbour.Y)] == floorTile {
						end = i + 1
					}else{
						start = i
					}
				}
			}
			if end < len(hall) {
				break
			}*/
			if e.tiles[int(tile.X)][int(tile.Y)] == floorTile {
				if geom.RectangleContainsInclusive(node.Right().Area(), tile) {
					end = i
					break
				}else{
					start = i + 1
				}
			}
		}
		
		hall = hall[start:end]
		for _, tile := range hall {
			e.tiles[int(tile.X)][int(tile.Y)] = floorTile
		}
	}
	return append(leftRooms, rightRooms...)
}

func generateEnvironment(e environment) environment {
	e.generateRooms(1)
	e.connectRooms(e.rooms.Root(), 0)
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