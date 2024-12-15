package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"time"

	aoc "github.com/shraddhaag/aoc/library"
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
var blockersCausingInfiniteLoop = []Coordinate{}
var checkedBlockers = []Coordinate{}

var points = []Coordinate{}

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
	loopBlockers := []Coordinate{}
	nextPosition := getNextPosition(*guardPosition, direction)
	walk(nextPosition, direction, grid, &loopBlockers)

	endTime := time.Since(startTime)

	printInfiniteBlockers(grid, blockersCausingInfiniteLoop)
	printGrid(grid)

	fmt.Printf("blocks causing infinte loop: %v\n", blockersCausingInfiniteLoop)
	fmt.Printf("blocks count causing infinte loop: %v\n", len(blockersCausingInfiniteLoop))

	fmt.Printf("total spaces visited: %d\n\n", total)
	fmt.Printf("executed in %d microseconds\n", endTime.Microseconds())

	other()

	fmt.Println("Points: ", points)
}

func printInfiniteBlockers(grid *[][]rune, blockers []Coordinate) {
	for _, c := range blockers {
		(*grid)[c.Y][c.X] = 'O'
	}
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
		c := Coordinate{
			X: 4,
			Y: 3,
		}

		if nextPosition == c {
			shouldLog = true
		}
		logf("blockers: %v\n", blockers)

		if (*grid)[nextPosition.Y][nextPosition.X] != 'X' && contains(checkedBlockers, nextPosition) == false {
			if isStuckInLoop(direction, nextPosition, nextPosition, &[]Coordinate{}) {
				blockersCausingInfiniteLoop = append(blockersCausingInfiniteLoop, nextPosition)
			}
		}
		shouldLog = false

		walk(nextPosition, direction, grid, loopBlockers)
	} else {
		direction = getNextDirection(direction)
		nextPosition = getNextPosition(guardPosition, direction)
		walk(nextPosition, direction, grid, loopBlockers)
	}
}

func isStuckInLoop(currentDirection int, initialBlockerCoord, blockerCoord Coordinate, loopBlockers *[]Coordinate) bool {
	logf("blocker coord: %v\n", blockerCoord)

	if initialBlockerCoord == blockerCoord {
		checkedBlockers = append(checkedBlockers, blockerCoord)
	}
	newDirection := getNextDirection(currentDirection)

	if blockerCoord != initialBlockerCoord {
		switch newDirection {
		case Up:
			if blockerCoord.X+1 == initialBlockerCoord.X && blockerCoord.Y > initialBlockerCoord.Y {
				logln("Hit inital blocker, found inifinte loop")
				return true
			}
		case Down:
			if blockerCoord.X-1 == initialBlockerCoord.X && blockerCoord.Y < initialBlockerCoord.Y {
				logln("Hit inital blocker, found inifinte loop")
				return true
			}
		case Right:
			logln("Checking for initial blocker right")
			logf("Intial blocker: %v\n", initialBlockerCoord)
			logf("blocker coord: %v\n", blockerCoord)
			if blockerCoord.Y+1 == initialBlockerCoord.Y && blockerCoord.X < initialBlockerCoord.X {
				logln("Hit inital blocker, found inifinte loop")
				return true
			}
		case Left:
			if blockerCoord.Y-1 == initialBlockerCoord.Y && blockerCoord.X > initialBlockerCoord.X {
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
				innerBlock := Coordinate{X: blockerCoord.X + 1, Y: blockerPosition}
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
				innerBlock := Coordinate{X: blockerCoord.X - 1, Y: blockerPosition}
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
				innerBlock := Coordinate{X: blockerPosition, Y: blockerCoord.Y + 1}
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
				innerBlock := Coordinate{X: blockerPosition, Y: blockerCoord.Y - 1}
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

	fmt.Printf("Vertical blockers: %v\n", verticalBlockers)
	fmt.Printf("Horizontal blockers: %v\n", horizontalBlockers)
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
	if r != '.' && r != 'X' && r != '^' {
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

// NOT MY SOLUTION

func other() {
	input := aoc.Get2DGrid(aoc.ReadFileLineByLine("../input.txt"))
	ans1, path := findPath(input, findStartingPoint(input))
	ans2 := findNewObstacleCount(input, findStartingPoint(input), path)

	fmt.Println("answer for part 1: ", ans1)
	fmt.Println("answer for part 2: ", ans2)
	fmt.Println("points: ", points)
}

func findStartingPoint(input [][]string) point {
	for i, row := range input {
		for j, char := range row {
			switch char {
			case "^":
				return point{j, i, up}
			case "<":
				return point{j, i, left}

			case ">":
				return point{j, i, right}

			case "v":
				return point{j, i, down}

			}
		}
	}
	return point{-1, -1, up}
}

type point struct {
	x         int
	y         int
	direction int
}

type coordinates struct {
	x int
	y int
}

const (
	up    = 0
	down  = 1
	right = 2
	left  = 3
)

func findPath(input [][]string, start point) (int, map[coordinates]int) {
	path := make(map[coordinates]int)
	count := 0
	current := start
	for {
		if _, ok := path[coordinates{current.x, current.y}]; !ok {
			count++
			path[coordinates{current.x, current.y}] = current.direction
		}

		isValid, newCurrent := findNextStep(input, current)
		if !isValid {
			return count, path
		}

		current = newCurrent
	}
	return count, path
}

// main thing to note: we enounter a loop whenever we come
// across the same coordinates + direction.
func isLoop(input [][]string, start point) bool {
	path := make(map[point]struct{})
	path2 := make(map[coordinates]struct{})
	current := start
	for {
		// update path
		if _, ok := path[current]; !ok {
			path[current] = struct{}{}
		} else {
			return true
		}

		if _, ok := path2[coordinates{current.x, current.y}]; !ok {
			path2[coordinates{current.x, current.y}] = struct{}{}
		}

		valid, newCurrent := findNextStep(input, current)
		if !valid {
			return false
		}

		current = newCurrent
	}
	return false
}

func findNextStep(input [][]string, current point) (bool, point) {
	valid, possibleNext := getNextStepWithDirectionPreserved(input, current)
	if !valid {
		return false, possibleNext
	}

	switch input[possibleNext.y][possibleNext.x] {
	case "#":
		// this is really subtle, consider the below case
		// where > represents your current position + direction:
		// ....#.....
		// ........>#
		// ........#.
		// When you turn right and step, you again encounter
		// an obstacle (a '#'). This is a valid case and you
		// can not exit the loop at this point.
		// Instead, you turn right once MORE, an effective turn of
		// 180 degrees this time, and then continue forward.
		return findNextStep(input, turn90(input, current))
	case ".":
		return true, possibleNext
	case "^":
		return true, possibleNext
	}
	return false, possibleNext
}

// for the guard to be stuck in a loop, the new obstacle has to
// be placed on the guard's existing path (ie path figured out
// in the first part).
// obstacle placed at any other place will not change the guard's path.
func findNewObstacleCount(input [][]string, start point, path map[coordinates]int) int {
	count := 0
	obstanceMap := make(map[coordinates]struct{})
	for step, _ := range path {
		if step.x == start.x && step.y == start.y {
			continue
		}
		if input[step.y][step.x] == "." {

			input[step.y][step.x] = "#"
			if isLoop(input, start) {
				if _, ok := obstanceMap[step]; !ok {
					count++
					points = append(points, Coordinate{X: step.x, Y: step.y})
					obstanceMap[step] = struct{}{}
				}
			}
			input[step.y][step.x] = "."
		}
	}
	return count
}

func getNextStepWithDirectionPreserved(input [][]string, current point) (bool, point) {
	switch current.direction {
	case up:
		current.y -= 1
	case down:
		current.y += 1
	case right:
		current.x += 1
	case left:
		current.x -= 1
	}
	if current.x < 0 || current.y < 0 || current.x >= len(input[0]) || current.y >= len(input) {
		return false, current
	}
	return true, current
}

func turn90(input [][]string, current point) point {
	switch current.direction {
	case up:
		current.direction = right
	case down:
		current.direction = left
	case right:
		current.direction = down
	case left:
		current.direction = up
	}
	return current
}
