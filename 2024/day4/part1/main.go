package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	x = 'X'
)

var maxGridDepth int
var maxGridWidth int

func main() {
	runeGrid, err := readFile("test.txt")
	if err != nil {
		fmt.Printf("Error reading input: %s", err.Error())
		os.Exit(1)
	}

	maxGridDepth = len(*runeGrid)
	maxGridWidth = len((*runeGrid)[0])

	for _, runeList := range *runeGrid {
		for _, rune := range runeList {
			fmt.Print(string(rune))
		}
		fmt.Println()
	}
}

func findXmasCount(runeGrid *[][]rune) int {
	total := 0
	for gridIndex, runeList := range *runeGrid {
		for lineIndex := range runeList {
			total += startSearch(gridIndex, lineIndex, runeGrid)
		}
	}

	return 0
}

func startSearch(gridIndex, lineIndex int, runeGrid *[][]rune) int {
	currentRune := (*runeGrid)[gridIndex][lineIndex]
	if currentRune != 'X' {
		return 0
	}

	return 0
}

func searchAround(gridIndex, lineIndex int, runeGrid *[][]rune, searchValue rune) bool {
	positions := make([]Coordinate, 8)

	left := Coordinate{gridIndex, lineIndex - 1}
	topLeft := Coordinate{gridIndex - 1, lineIndex - 1}
	top := Coordinate{gridIndex - 1, lineIndex}
	topRight := Coordinate{gridIndex - 1, lineIndex + 1}
	right := Coordinate{gridIndex, lineIndex + 1}
	bottomRight := Coordinate{gridIndex - 1, lineIndex + 1}
	bottom := Coordinate{gridIndex - 1, lineIndex}
	bottomLeft := Coordinate{gridIndex - 1, lineIndex - 1}

	positions[0] = left
	positions[1] = topLeft
	positions[2] = top
	positions[3] = topRight
	positions[4] = right
	positions[5] = bottomRight
	positions[6] = bottom
	positions[7] = bottomLeft

	for _, position := range positions {
		if inGrid(position) {
			charAtPosition := (*runeGrid)[position.X][position.Y]
			if charAtPosition == searchValue {
				if searchValue == 'S' {
					return true
				}

				return searchAround(position.X, position.Y, runeGrid, *getNextSearch(charAtPosition))
			}
		}
	}

	return false
}

func inGrid(coord Coordinate) bool {
	return coord.X >= 0 && coord.X <= maxGridDepth && coord.Y >= 0 && coord.Y <= maxGridWidth
}

func getNextSearch(r rune) *rune {
	switch r {
	case 'M':
		a := 'A'
		return &a
	case 'A':
		s := 'S'
		return &s
	}

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
