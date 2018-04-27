package main

import (
	"strings"
	"github.com/nsf/termbox-go"
)

const (
	xAlignLeft uint8 = iota
	xAlignCentre
	xAlignRight
)

const (
	yAlignBelow uint8 = iota
	yAlignCentre
	yAlignAbove
)

type label struct {
	text string
	textColor termbox.Attribute
	textHighlight termbox.Attribute
	location point
	xAlign uint8
	yAlign uint8
}

func initLabel(t string, tc, th termbox.Attribute, loc point, xa, ya uint8) label {
	l := label{
		text: t,
		textColor: tc,
		textHighlight: th,
		location: loc,
		xAlign: xa,
		yAlign: ya,
	}
	return l
}

func initLocationlessLabel(t string, tc, th termbox.Attribute) label {
	l := label{
		text: t,
		textColor: tc,
		textHighlight: th,
	}
	return l
}

func labelsEqual(a, b label) bool {
	return ((strings.Compare(a.text, b.text) == 0) && a.textColor == b.textColor && a.textHighlight == b.textHighlight && pointsEqual(a.location, b.location) && (a.xAlign == b.xAlign) && (a.yAlign == b.yAlign))
}

func alignLabels(bounds rectangle, xa, ya uint8, ls ...label) []label {
	x := bounds.corners[upperLeft].x + (bounds.corners[upperRight].x - bounds.corners[upperLeft].x) * (float64(xa) / 2)
	for i := 0; i < len(ls); i++ {
		ls[i].location = initPoint(x, bounds.corners[upperLeft].y + (bounds.corners[lowerLeft].y - bounds.corners[upperLeft].y) * (float64(i + 1) / float64(len(ls) + 1)))
		ls[i].xAlign = xa
		ls[i].yAlign = ya
	}
	return ls
}