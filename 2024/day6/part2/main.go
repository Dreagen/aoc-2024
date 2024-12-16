package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"time"
)

const (
	Up = iota
	Down
	Left
	Right
)

var shouldLog = false
var direction = Up
var total = 0
var visit = 0
var verticalBlockers map[int][]int
var horizontalBlockers map[int][]int
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

	makeBlockerMaps(grid, blockers)

	guardPosition := findGuardPosition(grid)
	loopBlockers := []CoordinateWithDirection{}
	nextPosition := getNextPosition(*guardPosition, direction)
	walk(nextPosition, direction, grid, &loopBlockers)

	endTime := time.Since(startTime)

	printInfiniteBlockers(grid, blockersCausingInfiniteLoop)
	printGrid(grid)

	sort.Slice(blockersCausingInfiniteLoop, func(i, j int) bool {
		if blockersCausingInfiniteLoop[i].X == blockersCausingInfiniteLoop[j].X {
			return blockersCausingInfiniteLoop[i].Y < blockersCausingInfiniteLoop[j].Y
		}
		return blockersCausingInfiniteLoop[i].X < blockersCausingInfiniteLoop[j].X
	})
	fmt.Printf("blocks causing infinte loop: %v\n", blockersCausingInfiniteLoop)
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

	c := Coordinate{
		X: 4,
		Y: 1,
	}

	nextPosition := getNextPosition(guardPosition, direction)
	if guardPosition.Coordinate == c {
		fmt.Printf("first: at 4, 1, next position: %v, dirtection: %d\n", nextPosition, direction)
	}
	if inGrid(nextPosition, grid) == false {
		if guardPosition.Coordinate == c {
			fmt.Println("Out of grid")
		}
		return
	}

	if nextPosition.Coordinate == c {
		fmt.Println("turning on logging")
		shouldLog = true
	}
	if isBlocked(nextPosition, grid) == false {
		logf("blockers: %v\n", blockers)

		if (*grid)[nextPosition.Y][nextPosition.X] != 'X' && contains(checkedBlockers, nextPosition.Coordinate) == false {
			if isStuckInLoop(direction, nextPosition, nextPosition, &[]CoordinateWithDirection{}) {
				blockersCausingInfiniteLoop = append(blockersCausingInfiniteLoop, nextPosition)
			}
		}
		shouldLog = false

		walk(nextPosition, direction, grid, loopBlockers)
	} else {
		direction = getNextDirection(direction)
		newNextPosition := getNextPosition(guardPosition, direction)
		if isBlocked(newNextPosition, grid) {
			direction = getNextDirection(direction)
			newNextPosition = getNextPosition(guardPosition, direction)
		}
		if (*grid)[newNextPosition.Y][newNextPosition.X] != 'X' && contains(checkedBlockers, nextPosition.Coordinate) == false {
			if isStuckInLoop(direction, newNextPosition, newNextPosition, &[]CoordinateWithDirection{}) {
				blockersCausingInfiniteLoop = append(blockersCausingInfiniteLoop, newNextPosition)
			}
		}
		if newNextPosition.Coordinate == c {
			fmt.Println("turning on logging")
			shouldLog = true
		}
		if guardPosition.Coordinate == c {
			fmt.Printf("blocked: at 12, 15, next position: %v, dirtection: %d\n", nextPosition, direction)
		}
		walk(newNextPosition, direction, grid, loopBlockers)
	}
}

func isStuckInLoop(currentDirection int, initialBlockerCoord, blockerCoord CoordinateWithDirection, loopBlockers *[]CoordinateWithDirection) bool {
	logf("blocker coord: %v\n", blockerCoord)

	newDirection := getNextDirection(currentDirection)

	if blockerCoord != initialBlockerCoord {
		switch newDirection {
		case Up:
			if newDirection == initialBlockerCoord.Direction && blockerCoord.X+1 == initialBlockerCoord.X && blockerCoord.Y > initialBlockerCoord.Y {
				logln("Hit inital blocker, found inifinte loop")
				return true
			}
		case Down:
			if newDirection == initialBlockerCoord.Direction && blockerCoord.X-1 == initialBlockerCoord.X && blockerCoord.Y < initialBlockerCoord.Y {
				logln("Hit inital blocker, found inifinte loop")
				return true
			}
		case Right:
			if newDirection == initialBlockerCoord.Direction && blockerCoord.Y+1 == initialBlockerCoord.Y && blockerCoord.X < initialBlockerCoord.X {
				logln("Hit inital blocker, found inifinte loop")
				return true
			}
		case Left:
			if newDirection == initialBlockerCoord.Direction && blockerCoord.Y-1 == initialBlockerCoord.Y && blockerCoord.X > initialBlockerCoord.X {
				logln("Hit inital blocker, found inifinte loop")
				return true
			}
		}
	}

	for _, blocker := range *loopBlockers {
		if blocker == blockerCoord {
			logf("Found blocker in loop blockers, blocker: %v, blocker in list: %v\n", blockerCoord, blocker)
			return true
		}
	}
	*loopBlockers = append(*loopBlockers, blockerCoord)

	switch newDirection {
	case Up:
		logln("Case Up")
		currentPosition := blockerCoord.X + 1
		rowToCheck := verticalBlockers[currentPosition]

		sorted := make([]int, len(rowToCheck))
		copy(sorted, rowToCheck)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i] > sorted[j]
		})

		for _, blockerPosition := range sorted {
			logf("Checking blocker position: %d agaisnt currentPosition: %d - needs to be smaller\n", blockerPosition, blockerCoord.Y)
			if blockerPosition < blockerCoord.Y {
				innerBlock := CoordinateWithDirection{Coordinate{X: blockerCoord.X + 1, Y: blockerPosition}, Up}
				logf("blocker position: %d, was smaller than current position: %d\n", blockerPosition, blockerCoord.Y)
				return isStuckInLoop(newDirection, initialBlockerCoord, innerBlock, loopBlockers)
			}
		}

		logln("returning false")
		return false
	case Down:
		logln("Case Down")
		currentPosition := blockerCoord.X - 1
		rowToCheck := verticalBlockers[currentPosition]
		for _, blockerPosition := range rowToCheck {
			logf("Checking blocker position: %d agaisnt currentPosition: %d - needs to be larger\n", blockerPosition, blockerCoord.Y)
			if blockerPosition > blockerCoord.Y {
				innerBlock := CoordinateWithDirection{Coordinate{X: blockerCoord.X - 1, Y: blockerPosition}, Down}
				logf("blocker position: %d, was larger than current position: %d\n", blockerPosition, blockerCoord.Y)
				return isStuckInLoop(newDirection, initialBlockerCoord, innerBlock, loopBlockers)
			}
		}

		logln("returning false")
		return false
	case Right:
		logln("Case Right")
		currentPosition := blockerCoord.Y + 1
		rowToCheck := horizontalBlockers[currentPosition]
		for _, blockerPosition := range rowToCheck {
			logf("Checking blocker position: %d agaisnt currentPosition: %d - needs to be larger\n", blockerPosition, blockerCoord.X)
			if blockerPosition > blockerCoord.X {
				innerBlock := CoordinateWithDirection{Coordinate{X: blockerPosition, Y: blockerCoord.Y + 1}, Right}
				logf("inner block: %d\n", innerBlock)
				logf("blocker position: %d, was larger than current position: %d\n", blockerPosition, blockerCoord.X)
				return isStuckInLoop(newDirection, initialBlockerCoord, innerBlock, loopBlockers)
			}
		}

		logln("returning false")
		return false
	case Left:
		logln("Case Left")
		currentPosition := blockerCoord.Y - 1
		rowToCheck := horizontalBlockers[currentPosition]

		sorted := make([]int, len(rowToCheck))
		copy(sorted, rowToCheck)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i] > sorted[j]
		})

		for _, blockerPosition := range sorted {
			logf("Checking blocker position: %d agaisnt currentPosition: %d - needs to be smaller\n", blockerPosition, blockerCoord.X)
			if blockerPosition < blockerCoord.X {
				innerBlock := CoordinateWithDirection{Coordinate{X: blockerPosition, Y: blockerCoord.Y - 1}, Left}
				logf("inner block: %d\n", innerBlock)
				logf("blocker position: %d, was smaller than current position: %d\n", blockerPosition, blockerCoord.X)
				return isStuckInLoop(newDirection, initialBlockerCoord, innerBlock, loopBlockers)
			}
		}

		logln("returning false")
		return false
	}

	panic("unknown direction")
}

func makeBlockerMaps(grid *[][]rune, blockers []Coordinate) {
	verticalBlockers = make(map[int][]int, len((*grid)[0]))
	horizontalBlockers = make(map[int][]int, len(*grid))

	for _, blocker := range blockers {
		verticalBlockers[blocker.X] = append(verticalBlockers[blocker.X], blocker.Y)
		horizontalBlockers[blocker.Y] = append(horizontalBlockers[blocker.Y], blocker.X)
	}

	for i := range verticalBlockers {
		sort.Slice(verticalBlockers[i], func(a, b int) bool {
			return verticalBlockers[i][a] < verticalBlockers[i][b]
		})
	}

	for i := range horizontalBlockers {
		sort.Slice(horizontalBlockers[i], func(a, b int) bool {
			return horizontalBlockers[i][a] < horizontalBlockers[i][b]
		})
	}

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
	inGrid := coord.X >= 0 && coord.X < len((*grid)[0]) && coord.Y >= 0 && coord.Y < len(*grid)
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

func logf(format string, a ...any) {
	if shouldLog {
		fmt.Printf(format, a...)
	}
}

func logln(format string) {
	if shouldLog {
		fmt.Println(format)
	}
}
