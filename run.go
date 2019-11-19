package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Galdoba/utils"
)

const (
	cellWLine       = "00000000"
	figureTypePawn  = " *|__|* "
	figureTypeQueen = " *||||* "
)

var coordRange []string
var checkersMAP map[int]figure
var globalError string

func coordsValid(coord string) bool {
	for i := range coordRange {
		if coord == coordRange[i] {
			return true
		}
	}
	return false
}

type figure interface {
	figureLines() []string
	figureCoords() string
	figurePlayer() int
	figureMaximumRange() int
	figureID() int
	figureSetCoords(string)
}

type pawn struct {
	player       int
	id           int
	maximumRange int
	fCoords      string
}

func newPawn(player int, id int) *pawn {
	p := &pawn{}
	p.player = player
	p.id = id
	p.maximumRange = 1
	checkersMAP[id] = p
	return p
}

func (p *pawn) figureID() int {
	return p.id
}

func (p *pawn) figurePlayer() int {
	return p.player
}

func (p *pawn) figureMaximumRange() int {
	return p.maximumRange
}

func (p *pawn) figureLines() []string {
	var lines []string
	switch p.player {
	case 1:
		lines = append(lines, " *+  +* ")
		lines = append(lines, " *++++* ")
	case 2:
		lines = append(lines, " *-  -* ")
		lines = append(lines, " *----* ")
	default:
		globalError = "Unknown player " + strconv.Itoa(p.player)
	}
	return lines
}

func (p *pawn) figureCoords() string {
	return p.fCoords
}

func (p *pawn) figureSetCoords(coords string) {
	if !coordsValid(coords) {
		globalError = coords + " is not Valid"
		return
	}
	p.fCoords = coords
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

func reverseSlice(sl []string) []string {
	for i := len(sl)/2 - 1; i >= 0; i-- {
		opp := len(sl) - 1 - i
		sl[i], sl[opp] = sl[opp], sl[i]
	}
	return sl
}

func pawnLines(player int) []string {
	var lines []string
	switch player {
	case 1:

		lines = append(lines, " *+  +* ")
		lines = append(lines, " *++++* ")

	case 2:
		lines = append(lines, " *-  -* ")
		lines = append(lines, " *----* ")
	default:
		globalError = "Unknown player " + strconv.Itoa(player)
	}
	return lines
}

func (b *board) update() {
	for i := range coordRange {
		lines := getBlankLines(b.cellMAP[coordRange[i]])
		//lines[0] = placeTag(lines[0], coordRange[i])
		b.cellMAP[coordRange[i]].lines = lines
		for _, v := range checkersMAP {

			if v.figureCoords() != coordRange[i] {

				continue
			}
			//pawnLines := pawnLines(v.player)
			pawnLines := v.figureLines()
			b.cellMAP[coordRange[i]].lines[1] = pawnLines[0]
			b.cellMAP[coordRange[i]].lines[2] = pawnLines[1]
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

func convertLetter(s string) int {
	s = strings.ToUpper(s)
	switch s {
	case "A":
		return 1
	case "B":
		return 2
	case "C":
		return 3
	case "D":
		return 4
	case "E":
		return 5
	case "F":
		return 6
	case "G":
		return 7
	case "H":
		return 8
	}
	return -1
}

func convertNumber(i int) string {
	switch i {
	case 1:
		return "A"
	case 2:
		return "B"
	case 3:
		return "C"
	case 4:
		return "D"
	case 5:
		return "E"
	case 6:
		return "F"
	case 7:
		return "G"
	case 8:
		return "H"
	}
	return "Z"
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
		fig := newPawn(pl, i)
		fig.figureSetCoords(startCoords[i])
	}
}

func clearTerm() {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	checkersMAP = make(map[int]figure)
	board := buildBoard()
	preapareGame()
	clearTerm()
	//fig := newpawn(1, 0)
	//fig.figureSetCoords("B1")
	board.update()
	drawBoard(board)
	//fig.figureSetCoords("C5")
	utils.InputString("тест ввода:")
	checkersMAP[9].figureSetCoords("E4")
	fmt.Println("////////////////////////")
	clearTerm()
	board.update()
	drawBoard(board)
	Move(checkersMAP[9], "E4")

	moveOptions(checkersMAP[9])
}

func isForward(f figure) int {
	if f.figurePlayer() == 1 {
		return 8
	}
	if f.figurePlayer() == 2 {
		return 1
	}
	return -1
}

func currentRow(f figure) int {
	coords := f.figureCoords()
	rowCol := strings.Split(coords, "")
	row, err := strconv.Atoi(rowCol[1])
	if err != nil {
		panic(err)
	}
	return row
}

//Move - двигает фигуру, если координаты валидны
func Move(f figure, newCoords string) {
	currentRow := currentRow(f)
	endRow := isForward(f)
	forwardRows := forwardRows(currentRow, endRow)
	fmt.Println(forwardRows)
}

func forwardRows(currentRow, endRow int) []int {
	var forwardRows []int
	if endRow == 8 {
		for currentRow < endRow {
			currentRow++
			forwardRows = append(forwardRows, currentRow)
		}
	}
	if endRow == 1 {
		for currentRow > endRow {
			currentRow++
			forwardRows = append(forwardRows, currentRow)
		}
	}

	return forwardRows
}

func moveOptions(f figure) []string {
	//var pCoords []string
	var allCoords []string
	for i := 1; i <= f.figureMaximumRange(); i++ {
		coords := f.figureCoords()
		row, col := coordsToRC(coords)
		allCoords = append(allCoords, rcToCoords(row+i, col+i))
		allCoords = append(allCoords, rcToCoords(row+i, col-i))
		allCoords = append(allCoords, rcToCoords(row-i, col-i))
		allCoords = append(allCoords, rcToCoords(row-i, col+i))
	}
	fmt.Println(allCoords)
	return allCoords
}

func coordsToRC(s string) (int, int) {
	rcStr := strings.Split(s, "")
	col := convertLetter(rcStr[0])
	row, _ := strconv.Atoi(rcStr[1])
	return row, col
}

func rcToCoords(row, col int) string {
	return convertNumber(col) + strconv.Itoa(row)
}

/*
движение:
1. Собираем ВСЕ клетки под угрозой
1а.
*/

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
