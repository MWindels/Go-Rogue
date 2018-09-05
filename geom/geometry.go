package geom

import (
	"fmt"
	"math"
)

const (
	upperLeft int = iota
	upperRight
	lowerRight
	lowerLeft
)

type Point struct {
	X float64
	Y float64
}

type Rectangle struct {
	corners [4]Point
}

func InitPoint(x, y float64) Point {
	p := Point{
		X: x,
		Y: y,
	}
	return p
}

func PointsEqual(a, b Point) bool {
	return (a.X == b.X && a.Y == b.Y)
}

func PointDistance(a, b Point) float64 {
	return math.Sqrt(math.Pow(a.X - b.X, 2) + math.Pow(a.Y - b.Y, 2))
}

func InitRectangle(x, y, width, height float64) Rectangle {
	r := Rectangle{
		corners: *new([4]Point),
	}
	r.corners[upperLeft] = InitPoint(x, y)
	r.corners[upperRight] = InitPoint(x + width, y)
	r.corners[lowerRight] = InitPoint(x + width, y + height)
	r.corners[lowerLeft] = InitPoint(x, y + height)
	return r
}

func (r Rectangle) UpperLeft() Point {
	return r.corners[upperLeft]
}

func (r Rectangle) UpperRight() Point {
	return r.corners[upperRight]
}

func (r Rectangle) LowerRight() Point {
	return r.corners[lowerRight]
}

func (r Rectangle) LowerLeft() Point {
	return r.corners[lowerLeft]
}

func (r Rectangle) Corner(i int) Point {
	if i < 0 || 4 <= i {
		panic(fmt.Sprintf("No corner with index %d.", i))
	}
	return r.corners[i]
}

func (r Rectangle) Width() float64 {
	return r.UpperRight().X - r.UpperLeft().X
}

func (r Rectangle) Height() float64 {
	return r.LowerLeft().Y - r.UpperLeft().Y
}

func RectanglesEqual(a, b Rectangle) bool {
	return PointsEqual(a.UpperLeft(), b.UpperLeft()) && PointsEqual(a.UpperRight(), b.UpperRight()) && PointsEqual(a.LowerRight(), b.LowerRight()) && PointsEqual(a.LowerLeft(), b.LowerLeft())
}

/*func DiscretizeRectangle(r Rectangle) Rectangle {
	return InitRectangle(math.Floor(r.UpperLeft().X), math.Floor(r.UpperLeft().Y), math.Ceil(r.UpperRight().X) - math.Floor(r.UpperLeft().X), math.Ceil(r.LowerLeft().Y) - math.Floor(r.UpperLeft().Y))
}*/

func ScaleRectangle(r Rectangle, xScale, yScale float64) Rectangle {
	return InitRectangle(xScale * r.UpperLeft().X, yScale * r.UpperLeft().Y, xScale * r.Width(), yScale * r.Height())
}

func RectangleContains(r Rectangle, p Point) bool {
	return (r.UpperLeft().X < p.X && p.X < r.UpperRight().X) && (r.UpperLeft().Y < p.Y && p.Y < r.LowerLeft().Y)
}

func RectangleContainsInclusive(r Rectangle, p Point) bool {
	return (r.UpperLeft().X <= p.X && p.X <= r.UpperRight().X) && (r.UpperLeft().Y <= p.Y && p.Y <= r.LowerLeft().Y)
}

func RectangleContainsLowerInclusive(r Rectangle, p Point) bool {
	return (r.UpperLeft().X <= p.X && p.X < r.UpperRight().X) && (r.UpperLeft().Y <= p.Y && p.Y < r.LowerLeft().Y)
}

func RectanglesOverlap(a, b Rectangle) bool {
	for i := 0; i < 4; i++ {
		if RectangleContains(a, b.Corner(i)) || RectangleContains(b, a.Corner(i)) {
			return true
		}
	}
	return false
}

func RectanglesOverlapInclusive(a, b Rectangle) bool {
	for i := 0; i < 4; i++ {
		if RectangleContainsInclusive(a, b.Corner(i)) || RectangleContainsInclusive(b, a.Corner(i)) {
			return true
		}
	}
	return false
}

func RectanglesIntersection(a, b Rectangle) Rectangle {
	upperLeft := InitPoint(math.Max(a.UpperLeft().X, b.UpperLeft().X), math.Max(a.UpperLeft().Y, b.UpperLeft().Y))
	lowerRight := InitPoint(math.Min(a.LowerRight().X, b.LowerRight().X), math.Min(a.LowerRight().Y, b.LowerRight().Y))
	return InitRectangle(upperLeft.X, upperLeft.Y, lowerRight.X - upperLeft.X, lowerRight.Y - upperLeft.Y)
}

func RectangleDistance(a, b Rectangle) float64 {
	if a.LowerRight().X < b.UpperLeft().X && a.LowerRight().Y < b.UpperLeft().Y {
		return PointDistance(a.LowerRight(), b.UpperLeft())		//above and left
	}else if a.LowerLeft().X > b.UpperRight().X && a.LowerLeft().Y < b.UpperRight().Y {
		return PointDistance(a.LowerLeft(), b.UpperRight())		//above and right
	}else if a.UpperLeft().X > b.LowerRight().X && a.UpperLeft().Y > b.LowerRight().Y {
		return PointDistance(a.UpperLeft(), b.LowerRight())		//below and right
	}else if a.UpperRight().X < b.LowerLeft().X && a.UpperRight().Y > b.LowerLeft().Y {
		return PointDistance(a.UpperRight(), b.LowerLeft())		//below and left
	}else if a.LowerLeft().Y < b.UpperLeft().Y {
		return b.UpperLeft().Y - a.LowerLeft().Y	//above
	}else if a.UpperLeft().X > b.UpperRight().X {
		return a.UpperLeft().X - b.UpperRight().X	//right
	}else if a.UpperLeft().Y > b.LowerLeft().Y {
		return a.UpperLeft().Y - b.LowerLeft().Y	//below
	}else if a.UpperRight().X < b.UpperLeft().X {
		return b.UpperLeft().X - a.UpperRight().X	//left
	}
	return 0
}

func NearestRectangles(r Rectangle, rs ...Rectangle) ([]Rectangle, float64) {
	distance := math.Inf(1)
	var nearest []Rectangle
	for _, current := range rs {
		currentDistance := RectangleDistance(r, current)
		if float32(currentDistance) < float32(distance) {
			distance = currentDistance
			nearest = []Rectangle{current}
		}else if float32(currentDistance) == float32(distance) {
			nearest = append(nearest, current)
		}
	}
	return nearest, distance
}