package main

const (
	upperLeft int = iota
	upperRight
	lowerRight
	lowerLeft
)

type point struct {
	x float64
	y float64
}

type rectangle struct {
	corners [4]point
}

func initPoint(x, y float64) point {
	p := point{
		x: x,
		y: y,
	}
	return p
}

func pointsEqual(a, b point) bool {
	return (a.x == b.x && a.y == b.y)
}

func initRectangle(x, y, width, height float64) rectangle {
	r := rectangle{
		corners: *new([4]point),
	}
	r.corners[upperLeft] = initPoint(x, y)
	r.corners[upperRight] = initPoint(x + width, y)
	r.corners[lowerRight] = initPoint(x + width, y + height)
	r.corners[lowerLeft] = initPoint(x, y + height)
	return r
}

func rectanglesEqual(a, b rectangle) bool {
	return (a.corners[upperLeft] == b.corners[upperLeft] && a.corners[upperRight] == b.corners[upperRight] && a.corners[lowerRight] == b.corners[lowerRight] && a.corners[lowerLeft] == b.corners[lowerLeft])
}

func scaleRectangle(r rectangle, xScale, yScale float64) rectangle {
	return initRectangle(xScale * r.corners[upperLeft].x, yScale * r.corners[upperLeft].y, xScale * (r.corners[upperRight].x - r.corners[upperLeft].x), yScale * (r.corners[lowerLeft].y - r.corners[upperLeft].y))
}

func rectangleContains(r rectangle, p point) bool {
	return ((r.corners[upperLeft].x < p.x && p.x < r.corners[upperRight].x) && (r.corners[upperLeft].y < p.y && p.y < r.corners[lowerLeft].y))
}