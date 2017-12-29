package main

import (
	"strings"
	"github.com/nsf/termbox-go"
)

const (
	xAlignLeft = iota
	xAlignCentre
	xAlignRight
)

const (
	yAlignBelow = iota
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

func labelsEqual(a, b label) bool {
	return ((strings.Compare(a.text, b.text) == 0) && a.textColor == b.textColor && a.textHighlight == b.textHighlight && pointsEqual(a.location, b.location) && (a.xAlign == b.xAlign) && (a.yAlign == b.yAlign))
}