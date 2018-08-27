package main

import (
	"strings"
	"github.com/nsf/termbox-go"
	"github.com/mwindels/go-rogue/geom"
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
	location geom.Point
	xAlign uint8
	yAlign uint8
}

func initLabel(t string, tc, th termbox.Attribute, loc geom.Point, xa, ya uint8) label {
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
	return ((strings.Compare(a.text, b.text) == 0) && a.textColor == b.textColor && a.textHighlight == b.textHighlight && geom.PointsEqual(a.location, b.location) && (a.xAlign == b.xAlign) && (a.yAlign == b.yAlign))
}

func alignLabels(bounds geom.Rectangle, xa, ya uint8, ls ...label) []label {
	x := bounds.UpperLeft().X + bounds.Width() * (float64(xa) / 2)
	for i := 0; i < len(ls); i++ {
		ls[i].location = geom.InitPoint(x, bounds.UpperLeft().Y + bounds.Height() * (float64(i + 1) / float64(len(ls) + 1)))
		ls[i].xAlign = xa
		ls[i].yAlign = ya
	}
	return ls
}