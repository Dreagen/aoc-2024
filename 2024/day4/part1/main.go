package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	x = 'X'
)

var maxGridHeight int
var maxGridWidth int

func main() {
	// runeGrid, err := readFile("test.txt")
	runeGrid, err := readFile("../input.txt")
	if err != nil {
		fmt.Printf("Error reading input: %s", err.Error())
		os.Exit(1)
	}

	maxGridHeight = len(*runeGrid)
	maxGridWidth = len((*runeGrid)[0])

	for _, runeList := range *runeGrid {
		for _, rune := range runeList {
			fmt.Print(string(rune))
		}
		fmt.Println()
	}

	fmt.Printf("max height: %d, max width: %d\n", maxGridHeight, maxGridWidth)

	total := findXmasCount(runeGrid)

	fmt.Println()
	fmt.Printf("Total xmas: %d\n", total)
}

func findXmasCount(runeGrid *[][]rune) int {
	total := 0
	for gridIndex, runeList := range *runeGrid {
		for lineIndex := range runeList {
			total += startSearch(gridIndex, lineIndex, runeGrid)
		}
	}

	return total
}

func startSearch(yIndex, xIndex int, runeGrid *[][]rune) int {
	currentRune := (*runeGrid)[yIndex][xIndex]
	if currentRune != 'X' {
		return 0
	}

	// fmt.Println("Found X")

	return searchAround(yIndex, xIndex, runeGrid, 'M', nil)
}

func searchAround(yIndex, xIndex int, runeGrid *[][]rune, searchValue rune, alreadyMovingInDirection *int) int {
	positions := make([]Coordinate, 8)

	left := Coordinate{yIndex, xIndex - 1}
	topLeft := Coordinate{yIndex - 1, xIndex - 1}
	top := Coordinate{yIndex - 1, xIndex}
	topRight := Coordinate{yIndex - 1, xIndex + 1}
	right := Coordinate{yIndex, xIndex + 1}
	bottomRight := Coordinate{yIndex + 1, xIndex + 1}
	bottom := Coordinate{yIndex + 1, xIndex}
	bottomLeft := Coordinate{yIndex + 1, xIndex - 1}

	positions[0] = left
	positions[1] = topLeft
	positions[2] = top
	positions[3] = topRight
	positions[4] = right
	positions[5] = bottomRight
	positions[6] = bottom
	positions[7] = bottomLeft

	var positionsToSearch []Coordinate
	if alreadyMovingInDirection != nil {
		positionsToSearch = append(positionsToSearch, positions[*alreadyMovingInDirection])
	} else {
		positionsToSearch = positions
	}

	total := 0
	for i, position := range positionsToSearch {
		var direction int
		if alreadyMovingInDirection != nil {
			direction = *alreadyMovingInDirection
		} else {
			direction = i
		}

		if inGrid(position) {
			charAtPosition := (*runeGrid)[position.X][position.Y]
			if charAtPosition == searchValue {
				// fmt.Printf("Found %c in position %s\n", charAtPosition, printDirection(direction))
				if searchValue == 'S' {
					fmt.Printf("Found XMAS following direction %s\n", printDirection(direction))
					total++
					continue
				}

				total += searchAround(position.X, position.Y, runeGrid, *getNextSearch(charAtPosition), &direction)
			}
		}
	}

	return total
}

func printDirection(i int) string {
	switch i {
	case 0:
		return "left"
	case 1:
		return "up left"
	case 2:
		return "up"
	case 3:
		return "up right"
	case 4:
		return "right"
	case 5:
		return "down right"
	case 6:
		return "down"
	case 7:
		return "down left"
	}

	return "invalid direction"
}

func inGrid(coord Coordinate) bool {
	inGrid := coord.X >= 0 && coord.X < maxGridHeight && coord.Y >= 0 && coord.Y < maxGridWidth
	return inGrid
}

func getNextSearch(r rune) *rune {
	switch r {
	case 'X':
		m := 'M'
		return &m
	case 'M':
		a := 'A'
		return &a
	case 'A':
		s := 'S'
		return &s
	case 'S':
		m := 'M'
		return &m
	}

	fmt.Printf("couldn't get next search for character: %c\n", r)

	return nil
}

type Coordinate struct {
	X int
	Y int
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
