// https://adventofcode.com/2024/day/4#part2
package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const (
	x = 'X'
	m = 'M'
	a = 'A'
	s = 'S'
)

var maxGridHeight int
var maxGridWidth int

func main() {
	startTime := time.Now()

	// runeGrid, err := readFile("../test.txt")
	runeGrid, err := readFile("../input.txt")
	if err != nil {
		fmt.Printf("Error reading input: %s", err.Error())
		os.Exit(1)
	}

	maxGridHeight = len(*runeGrid)
	maxGridWidth = len((*runeGrid)[0])

	total := findXmasCount(runeGrid)

	fmt.Println()
	fmt.Printf("Execution time: %d microseconds\n", time.Since(startTime).Microseconds())
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
	if currentRune != 'A' {
		return 0
	}

	return searchAround(yIndex, xIndex, runeGrid)
}

func searchAround(yIndex, xIndex int, runeGrid *[][]rune) int {
	positions := make([]Coordinate, 8)

	topLeft := Coordinate{yIndex - 1, xIndex - 1}
	topRight := Coordinate{yIndex - 1, xIndex + 1}
	bottomRight := Coordinate{yIndex + 1, xIndex + 1}
	bottomLeft := Coordinate{yIndex + 1, xIndex - 1}

	if inGrid(topLeft) == false || inGrid(topRight) == false || inGrid(bottomRight) == false || inGrid(bottomLeft) == false {
		return 0
	}

	positions[0] = topLeft
	positions[1] = topRight
	positions[2] = bottomRight
	positions[3] = bottomLeft

	charAtTopLeft := (*runeGrid)[positions[0].X][positions[0].Y]
	charAtTopRight := (*runeGrid)[positions[1].X][positions[1].Y]
	charAtBottomRight := (*runeGrid)[positions[2].X][positions[2].Y]
	charAtBottomLeft := (*runeGrid)[positions[3].X][positions[3].Y]

	if charAtTopLeft != m && charAtTopLeft != s {
		return 0
	}

	if charAtTopRight != m && charAtTopRight != s {
		return 0
	}

	if isCorrectOppositeCorner(charAtTopLeft, charAtBottomRight) == false {
		return 0
	}

	if isCorrectOppositeCorner(charAtTopRight, charAtBottomLeft) == false {
		return 0
	}

	return 1
}

func isCorrectOppositeCorner(r, r1 rune) bool {
	switch r {
	case m:
		return r1 == s
	case s:
		return r1 == m
	}

	fmt.Printf("unexpected char when checking opposite corner: %c\n", r)
	return false
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

func printGrid(runeGrid *[][]rune) {
	for _, runeList := range *runeGrid {
		for _, rune := range runeList {
			fmt.Print(string(rune))
		}
		fmt.Println()
	}
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
