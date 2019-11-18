package main

import (
	"fmt"
	"strconv"
)

const (
	cellWLine       = "00000000"
	figureTypePawn  = " *|__|* "
	figureTypeQueen = " *||||* "
)

var coordRange []string
var checkersMAP map[int]*figure
var globalError string

func coordsValid(coord string) bool {
	for i := range coordRange {
		if coord == coordRange[i] {
			return true
		}
	}
	return false
}

type figure struct {
	player  int
	id      int
	fType   string
	fCoords string
}

func newFigure(player int, id int) *figure {
	fgr := &figure{}
	fgr.player = player
	fgr.id = id
	fgr.fType = figureTypePawn
	checkersMAP[id] = fgr
	return fgr
}

func (fgr *figure) setCoords(coords string) {
	if !coordsValid(coords) {
		globalError = coords + " is not Valid"
		return
	}
	fgr.fCoords = coords
}

type board struct {
	cellMAP map[string]*cell
}

func buildBoard() *board {
	b := &board{}
	b.cellMAP = make(map[string]*cell)
	colr := -1
	cellColor := "B"
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			l := letterRange()[(r)]
			coordRange = append(coordRange, l+strconv.Itoa(c+1))
			if colr < 0 {
				cellColor = "W"
			} else {
				cellColor = "B"
			}
			coords := l + strconv.Itoa(c+1)
			cl := newCell(coords, cellColor)
			b.cellMAP[coords] = cl
			colr = -1 * colr
		}
		colr = -1 * colr
	}
	return b
}

func drawBoard(b *board) {
	lRange := letterRange()
	nRange := reverseSlice(numberRange())
	for row := 0; row < 8; row++ {
		for i := 0; i < 4; i++ {
			for col := 0; col < 8; col++ {
				coords := lRange[col] + nRange[row]
				fmt.Print(b.cellMAP[coords].lines[i])
			}
			fmt.Print("\n")
		}
	}
}

func pawnLines(fType string, player int) []string {
	var lines []string
	switch fType {
	case figureTypePawn:
		if player == 1 {
			lines = append(lines, " *+  +* ")
			lines = append(lines, " *++++* ")
		}
		if player == 2 {
			lines = append(lines, " *-  -* ")
			lines = append(lines, " *----* ")
		}
	case figureTypeQueen:
		if player == 1 {
			lines = append(lines, " *+||+* ")
			lines = append(lines, " *++++* ")
		}
		if player == 2 {
			lines = append(lines, " *-||-* ")
			lines = append(lines, " *----* ")
		}
	default:
		globalError = "Unknown type: " + fType + " or player " + strconv.Itoa(player)
	}
	return lines
}

func (b *board) update() {
	for i := range coordRange {
		lines := getBlankLines(b.cellMAP[coordRange[i]])
		lines[0] = placeTag(lines[0], coordRange[i])
		b.cellMAP[coordRange[i]].lines = lines
		for _, v := range checkersMAP {

			if v.fCoords != coordRange[i] {

				continue
			}

			b.cellMAP[coordRange[i]].lines[1] = v.fType
		}

	}
}

func placeTag(s string, tag string) string {
	bts := []byte(s)
	return "\\" + tag + "/" + string(bts[4]) + string(bts[5]) + string(bts[6]) + string(bts[7])
}

func letterRange() []string {
	lRange := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	return lRange
}

func numberRange() []string {
	lRange := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	return lRange
}

type cell struct {
	coords string
	color  string
	lines  []string
}

func newCell(coords, color string) *cell {
	cl := &cell{}
	cl.color = color
	cl.coords = coords
	cl.lines = []string{"1", "2", "3", "4"}
	cl.lines = getBlankLines(cl)
	// for i := 0; i < 4; i++ {
	// 	if cl.color == "W" {
	// 		cl.lines[i] = cellWLine
	// 	} else {
	// 		switch i {
	// 		case 0:
	// 			cl.lines[i] = "  " + cl.coords + "    "
	// 		case 1:
	// 			cl.lines[i] = "        "
	// 		case 2:
	// 			cl.lines[i] = "        "
	// 		case 3:
	// 			cl.lines[i] = "        "
	// 		}
	// 	}
	// }
	return cl
}

func getBlankLines(cl *cell) []string {
	var lines []string
	switch cl.color {
	case "B":
		lines = []string{"        ", "        ", "        ", "        "}
	case "W":
		lines = []string{"00000000", "00000000", "00000000", "00000000"}
	}
	return lines
}

func main() {
	checkersMAP = make(map[int]*figure)
	board := buildBoard()

	fig := newFigure(1, 0)
	fig.setCoords("B1")
	board.update()
	drawBoard(board)
}

func reverseSlice(sl []string) []string {
	for i := len(sl)/2 - 1; i >= 0; i-- {
		opp := len(sl) - 1 - i
		sl[i], sl[opp] = sl[opp], sl[i]
	}
	return sl
}

/*
oooooooo
        oooooooo
 *||||* oooooooo
 *++++* oooooooo
        oooooooo
oooooooo
        oooooooo
 *|  |* oooooooo
 *++++* oooooooo
        oooooooo


oooooooo
        oooooooo
 *-||-* oooooooo
 *----* oooooooo
        oooooooo
oooooooo
        oooooooo
 *-  -* oooooooo
 *----* oooooooo
        oooooooo
oooooooo


*/
