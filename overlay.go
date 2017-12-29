package main

import "github.com/nsf/termbox-go"

const (
	opaque = 1 << iota
	borderless
)

const (
	displayNothing = iota
	displayTitle
	displayEnvironment
)

type canvasConstants struct {
	labels []label
}

type canvas struct {
	attributes uint8
	borderCell termbox.Cell
	border rectangle
	variableContents uint8
	constantContents canvasConstants
}

type overlay struct {
	canvases []canvas
}

func initCanvasConstants() canvasConstants {
	cc := canvasConstants{
		labels: make([]label, 0, 1),
	}
	return cc
}

func addLabels(c canvasConstants, ls ...label) canvasConstants {
	c.labels = append(c.labels, ls...)
	return c
}

func canvasConstantsEqual(a, b canvasConstants) bool {
	if len(a.labels) != len(b.labels) {
		return false
	}
	for i := 0; i < len(a.labels); i++ {
		if !labelsEqual(a.labels[i], b.labels[i]) {
			return false
		}
	}
	return true
}

func initCanvas(a uint8, t rune, tfg, tbg termbox.Attribute, r rectangle, vc uint8, cc canvasConstants) canvas {
	c := canvas{
		attributes: a,
		borderCell: termbox.Cell{Ch: t, Fg: tfg, Bg: tbg},
		border: r,
		variableContents: vc,
		constantContents: cc,
	}
	return c
}

func canvasesEqual(a, b canvas) bool {
	return (a.attributes == b.attributes &&
			a.borderCell.Ch == b.borderCell.Ch &&
			a.borderCell.Fg == b.borderCell.Fg &&
			a.borderCell.Bg == b.borderCell.Bg &&
			rectanglesEqual(a.border, b.border) &&
			a.variableContents == b.variableContents &&
			canvasConstantsEqual(a.constantContents, b.constantContents))
}

func initOverlay() overlay {
	o := overlay{
		canvases: make([]canvas, 0, 1),
	}
	return o
}

func addToOverlay(o overlay, cs ...canvas) overlay {
	for i := 0; i < len(cs); i++ {
		canvasExists := false
		for j := 0; j < len(o.canvases); j++ {
			if canvasesEqual(cs[i], o.canvases[j]) {
				canvasExists = true
				break
			}
		}
		if !canvasExists {
			o.canvases = append(o.canvases, cs[i])
		}
	}
	return o
}

func canvasLayerOverlaps(o overlay, i int) bool {
	if i > len(o.canvases) {
		i = len(o.canvases)
	}
	for n := 0; n < i; n++ {
		for m := 0; m < 4; m++ {
			if rectangleContains(o.canvases[i].border, o.canvases[n].border.corners[m]) || rectangleContains(o.canvases[n].border, o.canvases[i].border.corners[m]) {
				return true
			}
		}
	}
	return false
}