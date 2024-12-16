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

func main() {
	startTime := time.Now()
	grid, err := readFile("../test.txt")
	// grid, err := readFile("../input.txt")
	if err != nil {
		fmt.Printf("Error reading input: %s", err.Error())
		os.Exit(1)
	}

	guardPosition := findGuardPosition(grid)
	walk(*guardPosition, direction, grid)

	endTime := time.Since(startTime)

	printGrid(grid)

	fmt.Printf("total spaces visited: %d\n\n", total)
	fmt.Printf("executed in %d microseconds\n", endTime.Microseconds())
}

func walk(guardPosition Coordinate, direction int, grid *[][]rune) {
	if (*grid)[guardPosition.Y][guardPosition.X] != 'X' {
		total++
	}
	(*grid)[guardPosition.Y][guardPosition.X] = 'X'

	nextPosition := getNextPosition(guardPosition, direction)
	if inGrid(nextPosition, grid) == false {
		return
	}

	if isBlocked(nextPosition, grid) == false {
		walk(nextPosition, direction, grid)
	} else {
		direction = getNextDirection(direction)
		nextPosition = getNextPosition(guardPosition, direction)
		walk(nextPosition, direction, grid)
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

	for scanner.Scan() {
		runesInLine = []rune(scanner.Text())
		runeGrid = append(runeGrid, runesInLine)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &runeGrid, nil
}
