package geom

import "fmt"

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

func RectanglesEqual(a, b Rectangle) bool {
	return PointsEqual(a.UpperLeft(), b.UpperLeft()) && PointsEqual(a.UpperRight(), b.UpperRight()) && PointsEqual(a.LowerRight(), b.LowerRight()) && PointsEqual(a.LowerLeft(), b.LowerLeft())
}

func ScaleRectangle(r Rectangle, xScale, yScale float64) Rectangle {
	return InitRectangle(xScale * r.UpperLeft().X, yScale * r.UpperLeft().Y, xScale * (r.UpperRight().X - r.UpperLeft().X), yScale * (r.LowerLeft().Y - r.UpperLeft().Y))
}

func RectangleContains(r Rectangle, p Point) bool {
	return (r.UpperLeft().X < p.X && p.X < r.UpperRight().X) && (r.UpperLeft().Y < p.Y && p.Y < r.LowerLeft().Y)
}