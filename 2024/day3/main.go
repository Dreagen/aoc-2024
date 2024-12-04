package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	m            = 'm'
	u            = 'u'
	l            = 'l'
	openBracket  = '('
	closeBracket = ')'
	noMatch      = 'x'
)

func main() {
	fmt.Println("day3")

	runes, err := readFile("test.txt")
	if err != nil {
		fmt.Printf("Error reading input: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println(runes)
}

func findMatch(runes *[]rune) int {
	previousMatch := noMatch
	inMatch := false
	for _, currentRune := range *runes {
		if currentRune == m {
			if inMatch {
				inMatch = false
				previousMatch = noMatch
				continue
			}

			inMatch = true
			previousMatch = m
			continue
		}

		if currentRune == u && inMatch && previousMatch == m {
			previousMatch = u
			continue
		}

		if currentRune == l && inMatch && previousMatch == u {
			previousMatch = l
			continue
		}

		if currentRune == openBracket && inMatch && previousMatch == l {
			previousMatch = openBracket
			continue
		}
	}

	return -1
}

func readFile(fileName string) (*[]rune, error) {
	var runes []rune
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)

	for scanner.Scan() {
		runes = append(runes, rune(scanner.Text()[0]))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &runes, nil
}
