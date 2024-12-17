package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const (
	Up = iota
	Down
	Left
	Right
)

var gridWidth = 0
var gridHeight = 0

var shouldLog = false
var direction = Up
var total = 0
var visit = 0
var verticalBlockers map[int][]bool
var horizontalBlockers map[int][]bool
var blockers = []Coordinate{}
var potentialLoop = []Coordinate{}
var blockersCausingInfiniteLoop = []CoordinateWithDirection{}
var checkedBlockers = []Coordinate{}

func main() {
	startTime := time.Now()
	// grid, err := readFile("../test.txt")
	grid, err := readFile("../input.txt")
	if err != nil {
		fmt.Printf("Error reading input: %s", err.Error())
		os.Exit(1)
	}

	gridWidth = len((*grid)[0])
	gridHeight = len(*grid)

	makeBlockerMaps(blockers)

	guardPosition := findGuardPosition(grid)
	loopBlockers := []CoordinateWithDirection{}
	nextPosition := getNextPosition(*guardPosition, direction)
	walk(nextPosition, direction, grid, &loopBlockers)

	endTime := time.Since(startTime)

	printInfiniteBlockers(grid, blockersCausingInfiniteLoop)

	printGrid(grid)

	fmt.Printf("blocks count causing infinte loop: %v\n", len(blockersCausingInfiniteLoop))

	fmt.Printf("total spaces visited: %d\n\n", total)
	fmt.Printf("executed in %d milliseconds\n", endTime.Milliseconds())
}

func printInfiniteBlockers(grid *[][]rune, blockers []CoordinateWithDirection) {
	for _, c := range blockers {
		(*grid)[c.Y][c.X] = 'O'
	}
}

func walk(guardPosition CoordinateWithDirection, direction int, grid *[][]rune, loopBlockers *[]CoordinateWithDirection) {
	if (*grid)[guardPosition.Y][guardPosition.X] != 'X' {
		total++
	}
	(*grid)[guardPosition.Y][guardPosition.X] = 'X'

	nextPosition := getNextPosition(guardPosition, direction)

	if inGrid(nextPosition, grid) == false {
		return
	}

	if isBlocked(nextPosition, grid) == false {

		if (*grid)[nextPosition.Y][nextPosition.X] != 'X' && contains(checkedBlockers, nextPosition.Coordinate) == false {
			addBlockerToMap(nextPosition)
			if isStuckInLoop(direction, nextPosition, nextPosition, &[]CoordinateWithDirection{}) {
				blockersCausingInfiniteLoop = append(blockersCausingInfiniteLoop, nextPosition)
			}
			removeBlockerFromMap(nextPosition)
		}

		walk(nextPosition, direction, grid, loopBlockers)
	} else {
		direction = getNextDirection(direction)
		newNextPosition := getNextPosition(guardPosition, direction)
		if isBlocked(newNextPosition, grid) {
			direction = getNextDirection(direction)
			newNextPosition = getNextPosition(guardPosition, direction)
		}
		if (*grid)[newNextPosition.Y][newNextPosition.X] != 'X' && contains(checkedBlockers, nextPosition.Coordinate) == false {
			addBlockerToMap(newNextPosition)
			if isStuckInLoop(direction, newNextPosition, newNextPosition, &[]CoordinateWithDirection{}) {
				blockersCausingInfiniteLoop = append(blockersCausingInfiniteLoop, newNextPosition)
			}
			removeBlockerFromMap(newNextPosition)
		}

		walk(newNextPosition, direction, grid, loopBlockers)
	}
}

func addBlockerToMap(c CoordinateWithDirection) {
	verticalBlockers[c.X][c.Y] = true
	horizontalBlockers[c.Y][c.X] = true
}

func removeBlockerFromMap(c CoordinateWithDirection) {
	verticalBlockers[c.X][c.Y] = false
	horizontalBlockers[c.Y][c.X] = false
}

func isStuckInLoop(currentDirection int, initialBlockerCoord, blockerCoord CoordinateWithDirection, loopBlockers *[]CoordinateWithDirection) bool {

	newDirection := getNextDirection(currentDirection)

	for _, blocker := range *loopBlockers {
		if blocker == blockerCoord {
			return true
		}
	}
	*loopBlockers = append(*loopBlockers, blockerCoord)

	switch newDirection {
	case Up:
		currentPosition := blockerCoord.X + 1
		rowToCheck := verticalBlockers[currentPosition]

		for i := len(rowToCheck) - 1; i >= 0; i-- {
			if i >= blockerCoord.Y {
				continue
			}

			if rowToCheck[i] == true {
				innerBlock := CoordinateWithDirection{Coordinate{X: blockerCoord.X + 1, Y: i}, Up}
				return isStuckInLoop(newDirection, initialBlockerCoord, innerBlock, loopBlockers)
			}
		}

		return false
	case Down:
		currentPosition := blockerCoord.X - 1
		rowToCheck := verticalBlockers[currentPosition]

		for i := 0; i < len(rowToCheck); i++ {
			if i <= blockerCoord.Y {
				continue
			}

			if rowToCheck[i] == true {
				innerBlock := CoordinateWithDirection{Coordinate{X: blockerCoord.X - 1, Y: i}, Down}
				return isStuckInLoop(newDirection, initialBlockerCoord, innerBlock, loopBlockers)
			}
		}

		return false
	case Right:
		currentPosition := blockerCoord.Y + 1
		rowToCheck := horizontalBlockers[currentPosition]

		for i := 0; i < len(rowToCheck); i++ {
			if i <= blockerCoord.X {
				continue
			}

			if rowToCheck[i] == true {
				innerBlock := CoordinateWithDirection{Coordinate{X: i, Y: blockerCoord.Y + 1}, Right}
				return isStuckInLoop(newDirection, initialBlockerCoord, innerBlock, loopBlockers)
			}
		}

		return false
	case Left:
		currentPosition := blockerCoord.Y - 1
		rowToCheck := horizontalBlockers[currentPosition]

		for i := len(rowToCheck) - 1; i >= 0; i-- {
			if i >= blockerCoord.X {
				continue
			}

			if rowToCheck[i] == true {
				innerBlock := CoordinateWithDirection{Coordinate{X: i, Y: blockerCoord.Y - 1}, Left}
				return isStuckInLoop(newDirection, initialBlockerCoord, innerBlock, loopBlockers)
			}
		}

		return false
	}

	panic("unknown direction")
}

func makeBlockerMaps(blockers []Coordinate) {
	verticalBlockers = make(map[int][]bool, gridWidth)
	horizontalBlockers = make(map[int][]bool, gridHeight)

	for i := 0; i < gridWidth; i++ {
		verticalBlockers[i] = make([]bool, gridHeight)
	}
	for i := 0; i < gridHeight; i++ {
		horizontalBlockers[i] = make([]bool, gridWidth)
	}

	for _, blocker := range blockers {
		verticalBlockers[blocker.X][blocker.Y] = true
		horizontalBlockers[blocker.Y][blocker.X] = true
	}

	// for i := range verticalBlockers {
	// 	sort.Slice(verticalBlockers[i], func(a, b int) bool {
	// 		return verticalBlockers[i][a] < verticalBlockers[i][b]
	// 	})
	// }
	//
	// for i := range horizontalBlockers {
	// 	sort.Slice(horizontalBlockers[i], func(a, b int) bool {
	// 		return horizontalBlockers[i][a] < horizontalBlockers[i][b]
	// 	})
	// }

	// fmt.Printf("Vertical blockers: %v\n", verticalBlockers)
	// fmt.Printf("Horizontal blockers: %v\n", horizontalBlockers)
}

func getNextPosition(guardPosition CoordinateWithDirection, direction int) CoordinateWithDirection {
	switch direction {
	case Up:
		return CoordinateWithDirection{
			Coordinate: Coordinate{
				X: guardPosition.X,
				Y: guardPosition.Y - 1,
			},
			Direction: Up,
		}
	case Down:
		return CoordinateWithDirection{
			Coordinate: Coordinate{
				X: guardPosition.X,
				Y: guardPosition.Y + 1,
			},
			Direction: Down,
		}
	case Left:
		return CoordinateWithDirection{
			Coordinate: Coordinate{
				X: guardPosition.X - 1,
				Y: guardPosition.Y,
			},
			Direction: Left,
		}
	case Right:
		return CoordinateWithDirection{
			Coordinate: Coordinate{
				X: guardPosition.X + 1,
				Y: guardPosition.Y,
			},
			Direction: Right,
		}
	}

	panic("direction not known")
}

func getNextDirection(currectDirection int) int {
	switch currectDirection {
	case Up:
		return Right
	case Right:
		return Down
	case Down:
		return Left
	case Left:
		return Up
	}

	panic("direction not known")
}

func isBlocked(position CoordinateWithDirection, grid *[][]rune) bool {
	r := (*grid)[position.Y][position.X]
	if r != '.' && r != 'X' && r != '^' {
		return true
	}

	return false
}

func inGrid(coord CoordinateWithDirection, grid *[][]rune) bool {
	inGrid := coord.X >= 0 && coord.X < gridWidth && coord.Y >= 0 && coord.Y < gridHeight
	return inGrid
}

func findGuardPosition(grid *[][]rune) *CoordinateWithDirection {
	for y, row := range *grid {
		for x, rune := range row {
			if rune == '^' {
				return &CoordinateWithDirection{
					Coordinate: Coordinate{
						X: x,
						Y: y,
					},
					Direction: Up,
				}
			}
		}
	}

	return nil
}

type Coordinate struct {
	X int
	Y int
}

type CoordinateWithDirection struct {
	Coordinate
	Direction int
}

func printGrid(runeGrid *[][]rune) {
	for _, runeList := range *runeGrid {
		for _, rune := range runeList {
			fmt.Print(string(rune))
		}
		fmt.Println()
	}

	fmt.Println()
}

func readFile(fileName string) (*[][]rune, error) {
	var runesInLine []rune
	var runeGrid [][]rune

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	i := 0
	for scanner.Scan() {
		runesInLine = []rune(scanner.Text())
		for j, r := range runesInLine {
			if r == '#' {
				blockers = append(blockers, Coordinate{X: j, Y: i})
			}
		}
		runeGrid = append(runeGrid, runesInLine)
		i++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &runeGrid, nil
}

func contains(slice []Coordinate, item Coordinate) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
