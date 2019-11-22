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
	cellWLine          = "00000000"
	figureTypePawn     = " *|__|* "
	figureTypeQueen    = " *||||* "
	directionNorthEast = "NE"
	directionSouthEast = "SE"
	directionSouthWest = "SW"
	directionNorthWest = "NW"
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
	figureValidMoveDirections() []string
}

/////////////////////////PAWN

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

func (p *pawn) figureValidMoveDirections() []string {
	var dir []string
	switch p.player {
	case 1:
		dir = append(dir, directionNorthEast)
		dir = append(dir, directionNorthWest)
	case 2:
		dir = append(dir, directionSouthEast)
		dir = append(dir, directionSouthWest)
	default:
		fmt.Println("Error: Unknown player detected ( player =", p.player, ")")
	}
	return dir
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

///////////////////KING

type king struct {
	player       int
	id           int
	maximumRange int
	fCoords      string
}

func newKing(player int, id int) *king {
	k := &king{}
	k.player = player
	k.id = id
	k.maximumRange = 7
	checkersMAP[id] = k
	return k
}

func (k *king) figureID() int {
	return k.id
}

func (k *king) figurePlayer() int {
	return k.player
}

func (k *king) figureMaximumRange() int {
	return k.maximumRange
}

func (k *king) figureLines() []string {
	var lines []string
	switch k.player {
	case 1:
		lines = append(lines, " *+||+* ")
		lines = append(lines, " *++++* ")
	case 2:
		lines = append(lines, " *-||-* ")
		lines = append(lines, " *----* ")
	default:
		globalError = "Unknown player " + strconv.Itoa(k.player)
	}
	return lines
}

func (k *king) figureValidMoveDirections() []string {
	return directionsAll()
}

func (k *king) figureCoords() string {
	return k.fCoords
}

func (k *king) figureSetCoords(coords string) {
	if !coordsValid(coords) {
		globalError = coords + " is not Valid"
		return
	}
	k.fCoords = coords
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
	promotePawns()
	for i := range coordRange {
		lines := getBlankLines(b.cellMAP[coordRange[i]])
		//lines[0] = placeTag(lines[0], coordRange[i])
		b.cellMAP[coordRange[i]].lines = lines
		for _, v := range checkersMAP {
			//fmt.Println(k)

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
	//seed := utils.RandomSeed()
	//fmt.Println(seed)
	checkersMAP = make(map[int]figure)
	board := buildBoard()
	//preapareGame()
	fig := newPawn(2, 0)
	fig.figureSetCoords("F7")
	fig2 := newPawn(1, 1)
	fig2.figureSetCoords("G6")
	fig3 := newPawn(1, 2)
	fig3.figureSetCoords("B7")
	clearTerm()
	board.update()
	drawBoard(board)

	utils.InputString("тест ввода:")
	//Move(figureByCoordinates("F7"), "E8")
	Attack(fig2, fig)
	board.update()
	drawBoard(board)
	Move(figureByCoordinates("B7"), "C8")
	//board.update()
	//clearTerm()
	board.update()
	drawBoard(board)
	readCommand("A2 take B2")
	board.update()
	drawBoard(board)

}

func readCommand(command string) (string, string, string) {
	command = strings.ToUpper(command)
	command = strings.TrimRight(command, " ")
	commPart := strings.Split(command, " ")
	return commPart[0], commPart[2], commPart[1]
}

func performAttack(command string) {
	//coord1, coord2, action := readCommand(command)

	fig1 := figureByCoordinates("D5")
	fig2 := figureByCoordinates("E4")
	pth := findPath(fig1.figureCoords(), fig2.figureCoords())
	start := fig1.figureCoords()
	dir := pth.direction
	dist := pth.distance
	if !testForAttack(fig1, pth) {
		fmt.Println("Attack impossible")
		return
	}
	//checkersMAP[fig2.figureID()] = nil
	delete(checkersMAP, fig2.figureID())
	checkersMAP[fig1.figureID()].figureSetCoords(navigateFrom(start, dir, dist+1))

}

func promotePawns() {
	for k, v := range checkersMAP {
		if _, ok := v.(*pawn); ok {
			fmt.Println("pawn", ok)
			continue
		}
		fmt.Println(k)
		coords := v.figureCoords()
		r, _ := coordsToRC(coords)
		pl := v.figurePlayer()
		lastRow := 0
		switch pl {
		case 1:
			lastRow = 8
		case 2:
			lastRow = 1
		}

		if r != lastRow {
			continue
		}
		id := v.figurePlayer()
		delete(checkersMAP, id)
		king := newKing(pl, id)
		king.figureSetCoords(coords)

		fmt.Println("row", r)
	}
}

func promotePawn(p pawn) {
	coords := p.figureCoords()
	row, _ := coordsToRC(coords)
	pl := p.figurePlayer()
	lastRow := 0
	switch pl {
	case 1:
		lastRow = 8
	case 2:
		lastRow = 1
	}
	if row != lastRow {
		return
	}
	id := p.figureID()
	delete(checkersMAP, id) //не сторого обязательно
	king := newKing(pl, id)
	king.figureSetCoords(coords)

}

func Attack(fig1, fig2 figure) {
	pth := findPath(fig1.figureCoords(), fig2.figureCoords())
	start := fig1.figureCoords()
	dir := pth.direction
	dist := pth.distance
	if !testForAttack(fig1, pth) {
		fmt.Println("Attack impossible")
		return
	}
	//checkersMAP[fig2.figureID()] = nil
	delete(checkersMAP, fig2.figureID())
	landCoords := navigateFrom(start, dir, dist+1)
	fig1.figureSetCoords(landCoords)
	fmt.Println("Land to", landCoords)
	fmt.Println("attack promotion start")
	if pwn, ok := fig1.(*pawn); ok {
		fmt.Println("attack promotion - go")
		promotePawn(*pwn)
		fmt.Println("attack promotion - end")
	}
}

func Move(fig figure, newCoords string) {
	fig.figureSetCoords(newCoords)
	//promotePawns()
	if pwn, ok := fig.(*pawn); ok {
		promotePawn(*pwn)
	}
}

type path struct {
	start     string
	direction string
	distance  int
}

func (p *path) isFree() bool {
	for i := 1; i <= p.distance; i++ {
		checkCoords := navigateFrom(p.start, p.direction, i)
		if !coordsValid(checkCoords) {
			return false
		}
		if !tileIsFree(checkCoords) {
			return false
		}
	}

	return true
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

func tileIsFree(s string) bool {
	if !coordsValid(s) {
		return false
	}
	for _, v := range checkersMAP {
		if v.figureCoords() != s {
			continue
		}
		return false
	}
	return true
}

func navigateFrom(currentCoords string, direction string, distance int) string {
	row, col := coordsToRC(currentCoords)
	addRow, addCol := directionToRC(direction)
	row = row + (addRow * distance)
	col = col + (addCol * distance)
	newCoords := rcToCoords(row, col)
	return newCoords
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

func attackList(f figure) []string {
	var attackCoords []string
	directions := directionsAll()
	maxRange := f.figureMaximumRange()
	for d := range directions {
		coords := f.figureCoords()
		for r := 0; r <= maxRange; r++ {
			newCoords := navigateFrom(coords, directions[d], r)
			if newCoords == f.figureCoords() {
				continue
			}
			if !coordsValid(newCoords) {
				continue
			}
			p := path{coords, directions[d], r}
			if p.isFree() {
				continue
			}
			fmt.Println("Path to", newCoords, "is not free")
			fmt.Println("Can attack", testForAttack(f, p))
			if testForAttack(f, p) {
				attackCoords = append(attackCoords, newCoords)
			}
			r = r + maxRange + 1

		}
	}
	fmt.Println("return attackCoords:", attackCoords)
	return attackCoords
}

func moveList(f figure) []string {
	var moveCoords []string
	directions := f.figureValidMoveDirections()
	maxRange := f.figureMaximumRange()
	for d := range directions {
		coords := f.figureCoords()
		for r := 0; r <= maxRange; r++ {
			newCoords := navigateFrom(coords, directions[d], r)
			if newCoords == f.figureCoords() {
				continue
			}
			if !coordsValid(newCoords) {
				continue
			}
			p := path{coords, directions[d], r}
			if p.isFree() {
				//fmt.Println("Path to", newCoords, "is free")
				moveCoords = append(moveCoords, newCoords)
			} else {
				// fmt.Println("Path to", newCoords, "is not free")
				// fmt.Println("Can attack", testForAttack(f, p))
				// if testForAttack(f, p) {
				// 	//	attackCoords = append(attackCoords, newCoords)
				// }
				r = r + maxRange + 1
			}
		}
	}

	fmt.Println("return moveCoords:", moveCoords)
	return moveCoords
}

// func projectMovement(f figure, directions ...string) []string {
// 	var moveCoords []string
// 	var attackCoords []string
// 	maxRange := f.figureMaximumRange()
// 	for d := range directions {
// 		coords := f.figureCoords()
// 		for r := 0; r <= maxRange; r++ {
// 			newCoords := navigateFrom(coords, directions[d], r)
// 			if newCoords == f.figureCoords() {
// 				continue
// 			}
// 			if !coordsValid(newCoords) {
// 				continue
// 			}
// 			p := path{coords, directions[d], r}
// 			if p.isFree() {
// 				fmt.Println("Path to", newCoords, "is free")
// 				moveCoords = append(moveCoords, newCoords)
// 			} else {
// 				fmt.Println("Path to", newCoords, "is not free")
// 				fmt.Println("Can attack", testForAttack(f, p))
// 				if testForAttack(f, p) {
// 					attackCoords = append(attackCoords, newCoords)
// 				}
// 				r = r + maxRange + 1
// 			}
// 		}
// 	}
// 	if len(attackCoords) > 0 {
// 		fmt.Println("return attackCoords:", attackCoords)
// 		return attackCoords
// 	}
// 	fmt.Println("return moveCoords:", moveCoords)
// 	return moveCoords
// }

//func attackList()

func testForAttack(f figure, p path) bool {
	attackCoords := navigateFrom(f.figureCoords(), p.direction, p.distance)
	if allied(f, figureByCoordinates(attackCoords)) {
		return false
	}
	landCoord := navigateFrom(f.figureCoords(), p.direction, p.distance+1)
	if !tileIsFree(landCoord) {
		return false
	}
	return true
}

func allied(f1 figure, f2 figure) bool {
	if f1.figurePlayer() == f2.figurePlayer() {
		return true
	}
	return false
}

func figureByCoordinates(coords string) figure {
	for _, v := range checkersMAP {
		if v.figureCoords() == coords {
			return v
		}
	}
	return nil

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
	allCoords = excludeInvalidCoords(allCoords)
	var freeCoords []string
	for i := range allCoords {
		if tileIsFree(allCoords[i]) {
			freeCoords = append(freeCoords, allCoords[i])
		}
	}
	return freeCoords
}

func coordsToRC(s string) (int, int) {
	rcStr := strings.Split(s, "")
	fmt.Println(rcStr)
	col := convertLetter(rcStr[0])
	row, _ := strconv.Atoi(rcStr[1])
	return row, col
}

func rcToCoords(row, col int) string {
	return convertNumber(col) + strconv.Itoa(row)
}

func directionToRC(dir string) (int, int) {
	addRow, addCol := 0, 0
	switch dir {
	case directionNorthEast:
		return 1, 1
	case directionSouthEast:
		return -1, 1
	case directionSouthWest:
		return -1, -1
	case directionNorthWest:
		return 1, -1
	}
	return addRow, addCol
}

func findPath(start string, end string) path { //TODO: функцию явно надо доработать - брут форс не наш метод!
	row, _ := coordsToRC(start)
	distance := 0
	direction := "No Path"
	var checkEnd []string
	for i := 1; i < 8; i++ {
		checkEnd = append(checkEnd, navigateFrom(start, directionNorthEast, i))
		checkEnd = append(checkEnd, navigateFrom(start, directionSouthEast, i))
		checkEnd = append(checkEnd, navigateFrom(start, directionSouthWest, i))
		checkEnd = append(checkEnd, navigateFrom(start, directionNorthWest, i))
	}
	checkEnd = excludeInvalidCoords(checkEnd)
	for j := range checkEnd {
		if checkEnd[j] != end {
			continue
		}
		endRow, _ := coordsToRC(end)
		distance = endRow - row
		if distance < 0 {
			distance = distance * -1
		}
		if navigateFrom(start, directionNorthEast, distance) == end {
			direction = directionNorthEast
		}
		if navigateFrom(start, directionSouthEast, distance) == end {
			direction = directionSouthEast
		}
		if navigateFrom(start, directionSouthWest, distance) == end {
			direction = directionSouthWest
		}
		if navigateFrom(start, directionNorthWest, distance) == end {
			direction = directionNorthWest
		}

	}
	p := path{start, direction, distance}
	return p
}

// func getPath(start string, end string) (string, int) { //TODO: функцию явно надо доработать - брут форс не наш метод!
// 	row, _ := coordsToRC(start)
// 	distance := 0
// 	direction := "No Path"
// 	var checkEnd []string
// 	for i := 1; i < 8; i++ {
// 		checkEnd = append(checkEnd, navigateFrom(start, directionNorthEast, i))
// 		checkEnd = append(checkEnd, navigateFrom(start, directionSouthEast, i))
// 		checkEnd = append(checkEnd, navigateFrom(start, directionSouthWest, i))
// 		checkEnd = append(checkEnd, navigateFrom(start, directionNorthWest, i))
// 	}
// 	checkEnd = excludeInvalidCoords(checkEnd)
// 	for j := range checkEnd {
// 		if checkEnd[j] != end {
// 			continue
// 		}
// 		endRow, _ := coordsToRC(end)
// 		distance = endRow - row
// 		if distance < 0 {
// 			distance = distance * -1
// 		}
// 		if navigateFrom(start, directionNorthEast, distance) == end {
// 			direction = directionNorthEast
// 		}
// 		if navigateFrom(start, directionSouthEast, distance) == end {
// 			direction = directionSouthEast
// 		}
// 		if navigateFrom(start, directionSouthWest, distance) == end {
// 			direction = directionSouthWest
// 		}
// 		if navigateFrom(start, directionNorthWest, distance) == end {
// 			direction = directionNorthWest
// 		}

// 	}
// 	return direction, distance
// }

func excludeInvalidCoords(coordsList []string) []string {
	var validCoords []string
	for i := range coordsList {
		if coordsValid(coordsList[i]) {
			validCoords = append(validCoords, coordsList[i])
		}
	}
	return validCoords
}

func directionsAll() []string {
	return []string{directionNorthEast, directionNorthWest, directionSouthEast, directionSouthWest}
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
