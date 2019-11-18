package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/Galdoba/utils"
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
	maxRow  int
	maxCol  int
}

func buildBoard() *board {
	b := &board{}
	b.cellMAP = make(map[string]*cell)
	b.maxRow = 8
	b.maxCol = 8
	colr := -1
	cellColor := "B"
	for r := 0; r < b.maxRow; r++ {
		for c := 0; c < b.maxCol; c++ {
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
	fmt.Println("  ++--------------------------------------------------------------------++")
	fmt.Println("  ++--------------------------------------------------------------------++")
	for row := 0; row < 8; row++ {
		for i := 0; i < 4; i++ {
			for col := 0; col < 8; col++ {
				if col == 0 {
					fmt.Print("  ||  ")
				}
				coords := lRange[col] + nRange[row]
				fmt.Print(b.cellMAP[coords].lines[i])
				if col == 7 {
					fmt.Print("  ||")
				}
				if i == 1 && col == 7 {
					fmt.Print(" " + nRange[row])
				}
			}
			fmt.Print("\n")
		}
	}
	fmt.Println("  ++--------------------------------------------------------------------++")
	fmt.Println("  ++--------------------------------------------------------------------++")
	fmt.Println("        A       B       C       D       E       F       G       H")
}

func figureLines(fType string, player int) []string {
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
		//lines[0] = placeTag(lines[0], coordRange[i])
		b.cellMAP[coordRange[i]].lines = lines
		for _, v := range checkersMAP {

			if v.fCoords != coordRange[i] {

				continue
			}
			figureLines := figureLines(v.fType, v.player)
			b.cellMAP[coordRange[i]].lines[1] = figureLines[0]
			b.cellMAP[coordRange[i]].lines[2] = figureLines[1]
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

func preapareGame() {
	startCoords := []string{"B1", "D1", "F1", "H1", "A2", "C2", "E2", "G2", "B3", "D3", "F3", "H3", //белые
		"A6", "C6", "E6", "G6", "B7", "D7", "F7", "H7", "A8", "C8", "E8", "G8"} //Черные
	pl := 1
	for i := range startCoords {
		if i >= 12 {
			pl = 2
		}
		fig := newFigure(pl, i)
		fig.setCoords(startCoords[i])
	}
}

func clearTerm() {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	checkersMAP = make(map[int]*figure)
	board := buildBoard()
	preapareGame()
	clearTerm()
	//fig := newFigure(1, 0)
	//fig.setCoords("B1")
	board.update()
	drawBoard(board)
	//fig.setCoords("C5")
	utils.InputString("тест ввода:")
	checkersMAP[9].setCoords("E4")
	fmt.Println("////////////////////////")
	clearTerm()
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
