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

var direction = Up
var total = 0
var visit = 0
var verticalBlockers map[int][]int
var horizontalBlockers map[int][]int
var blockers = []Coordinate{}
var potentialLoop = []Coordinate{}
var blockersCausingInfiniteLoop = []Coordinate{}

func main() {
	startTime := time.Now()
	grid, err := readFile("../test.txt")
	// grid, err := readFile("../input.txt")
	if err != nil {
		fmt.Printf("Error reading input: %s", err.Error())
		os.Exit(1)
	}

	makeBlockerMaps(grid, blockers)

	guardPosition := findGuardPosition(grid)
	loopBlockers := []Coordinate{}
	walk(*guardPosition, direction, grid, &loopBlockers)

	endTime := time.Since(startTime)

	printGrid(grid)

	fmt.Printf("blocks causing infinte loop: %v\n", blockersCausingInfiniteLoop)

	fmt.Printf("total spaces visited: %d\n\n", total)
	fmt.Printf("executed in %d microseconds\n", endTime.Microseconds())
}

func walk(guardPosition Coordinate, direction int, grid *[][]rune, loopBlockers *[]Coordinate) {
	if (*grid)[guardPosition.Y][guardPosition.X] != 'X' {
		total++
	}
	(*grid)[guardPosition.Y][guardPosition.X] = 'X'

	nextPosition := getNextPosition(guardPosition, direction)
	if inGrid(nextPosition, grid) == false {
		return
	}

	if isBlocked(nextPosition, grid) == false {
		if isStuckInLoop(direction, nextPosition, nextPosition, &[]Coordinate{}) {
			blockersCausingInfiniteLoop = append(blockersCausingInfiniteLoop, nextPosition)
		}

		walk(nextPosition, direction, grid, loopBlockers)
	} else {
		direction = getNextDirection(direction)
		nextPosition = getNextPosition(guardPosition, direction)
		walk(nextPosition, direction, grid, loopBlockers)
	}
}

func isStuckInLoop(currentDirection int, initialBlockerCoord, blockerCoord Coordinate, loopBlockers *[]Coordinate) bool {
	fmt.Printf("blocker coord: %v\n", blockerCoord)

	newDirection := getNextDirection(currentDirection)

	if blockerCoord != initialBlockerCoord {
		switch newDirection {
		case Up:
			if blockerCoord.X+1 == initialBlockerCoord.X && blockerCoord.Y > initialBlockerCoord.Y {
				fmt.Println("Hit inital blocker, found inifinte loop")
				return true
			}
		case Down:
			if blockerCoord.X-1 == initialBlockerCoord.X && blockerCoord.Y < initialBlockerCoord.Y {
				fmt.Println("Hit inital blocker, found inifinte loop")
				return true
			}
		case Right:
			if blockerCoord.Y+1 == initialBlockerCoord.Y && blockerCoord.X < initialBlockerCoord.X {
				fmt.Println("Hit inital blocker, found inifinte loop")
				return true
			}
		case Left:
			if blockerCoord.Y-1 == initialBlockerCoord.Y && blockerCoord.X > initialBlockerCoord.X {
				fmt.Println("Hit inital blocker, found inifinte loop")
				return true
			}
		}
	}

	for _, blocker := range *loopBlockers {
		if blocker == blockerCoord {
			fmt.Printf("Found blocker in loop blockers, blocker: %v, blocker in list: %v\n", blockerCoord, blocker)
			return true
		}
	}
	*loopBlockers = append(*loopBlockers, blockerCoord)

	switch newDirection {
	case Up:
		fmt.Println("Case Up")
		currentPosition := blockerCoord.X + 1
		rowToCheck := verticalBlockers[currentPosition]
		for _, blockerPosition := range rowToCheck {
			fmt.Printf("Checking blocker position: %d agaisnt currentPosition: %d - needs to be smaller\n", blockerPosition, blockerCoord.Y)
			if blockerPosition < blockerCoord.Y {
				innerBlock := Coordinate{X: blockerCoord.X + 1, Y: blockerPosition}
				fmt.Printf("blocker position: %d, was smaller than current position: %d\n", blockerPosition, blockerCoord.Y)
				return isStuckInLoop(newDirection, blockerCoord, innerBlock, loopBlockers)
			}
		}

		fmt.Println("returning false")
		return false
	case Down:
		fmt.Println("Case Down")
		currentPosition := blockerCoord.X - 1
		rowToCheck := verticalBlockers[currentPosition]
		for _, blockerPosition := range rowToCheck {
			fmt.Printf("Checking blocker position: %d agaisnt currentPosition: %d - needs to be larger\n", blockerPosition, blockerCoord.Y)
			if blockerPosition > blockerCoord.Y {
				innerBlock := Coordinate{X: blockerCoord.X - 1, Y: blockerPosition}
				fmt.Printf("blocker position: %d, was larger than current position: %d\n", blockerPosition, blockerCoord.Y)
				return isStuckInLoop(newDirection, blockerCoord, innerBlock, loopBlockers)
			}
		}

		fmt.Println("returning false")
		return false
	case Right:
		fmt.Println("Case Right")
		currentPosition := blockerCoord.Y + 1
		rowToCheck := horizontalBlockers[currentPosition]
		for _, blockerPosition := range rowToCheck {
			fmt.Printf("Checking blocker position: %d agaisnt currentPosition: %d - needs to be larger\n", blockerPosition, blockerCoord.X)
			if blockerPosition > blockerCoord.X {
				innerBlock := Coordinate{X: blockerPosition, Y: blockerCoord.Y + 1}
				fmt.Printf("inner block: %d\n", innerBlock)
				fmt.Printf("blocker position: %d, was larger than current position: %d\n", blockerPosition, blockerCoord.X)
				return isStuckInLoop(newDirection, blockerCoord, innerBlock, loopBlockers)
			}
		}

		fmt.Println("returning false")
		return false
	case Left:
		fmt.Println("Case Left")
		currentPosition := blockerCoord.Y - 1
		rowToCheck := horizontalBlockers[currentPosition]
		for _, blockerPosition := range rowToCheck {
			fmt.Printf("Checking blocker position: %d agaisnt currentPosition: %d - needs to be smaller\n", blockerPosition, blockerCoord.X)
			if blockerPosition < blockerCoord.X {
				innerBlock := Coordinate{X: blockerPosition, Y: blockerCoord.Y - 1}
				fmt.Printf("inner block: %d\n", innerBlock)
				fmt.Printf("blocker position: %d, was smaller than current position: %d\n", blockerPosition, blockerCoord.X)
				return isStuckInLoop(newDirection, blockerCoord, innerBlock, loopBlockers)
			}
		}

		fmt.Println("returning false")
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
}

func getNextPosition(guardPosition Coordinate, direction int) Coordinate {
	switch direction {
	case Up:
		return Coordinate{
			X: guardPosition.X,
			Y: guardPosition.Y - 1,
		}
	case Down:
		return Coordinate{
			X: guardPosition.X,
			Y: guardPosition.Y + 1,
		}
	case Left:
		return Coordinate{
			X: guardPosition.X - 1,
			Y: guardPosition.Y,
		}
	case Right:
		return Coordinate{
			X: guardPosition.X + 1,
			Y: guardPosition.Y,
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

func isBlocked(position Coordinate, grid *[][]rune) bool {
	r := (*grid)[position.Y][position.X]
	if r != '.' && r != 'X' {
		return true
	}

	return false
}

func inGrid(coord Coordinate, grid *[][]rune) bool {
	inGrid := coord.X >= 0 && coord.X < len((*grid)[0]) && coord.Y >= 0 && coord.Y < len(*grid)
	return inGrid
}

func findGuardPosition(grid *[][]rune) *Coordinate {
	for y, row := range *grid {
		for x, rune := range row {
			if rune == '^' {
				return &Coordinate{
					X: x,
					Y: y,
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
