package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const (
	m            = 'm'
	u            = 'u'
	l            = 'l'
	openBracket  = '('
	closeBracket = ')'
	comma        = ','
	noMatch      = 'x'
)

func main() {
	fmt.Println("day3")

	startTime := time.Now()

	runes, err := readFile("../input.txt")
	if err != nil {
		fmt.Printf("Error reading input: %s", err.Error())
		os.Exit(1)
	}

	result := findMatch(runes)
	endTime := time.Since(startTime)

	fmt.Printf("Result: %d\n", result)
	fmt.Printf("Took %d microseconds\n", endTime)
}

func findMatch(runes *[]rune) int {
	total := 0

	previousMatch := noMatch
	inMatch := false
	firstNumber := 0
	secondNumber := 0
	firstNumberMatched := false
	numbersMatchedInARow := 0

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

		if isInt(currentRune) && numbersMatchedInARow < 3 && inMatch && (previousMatch == openBracket || isInt(previousMatch)) {
			previousMatch = currentRune
			if firstNumberMatched == false {
				firstNumber = (firstNumber * 10) + runeToInt(currentRune)
				numbersMatchedInARow++
				continue
			}

			secondNumber = (secondNumber * 10) + runeToInt(currentRune)
			numbersMatchedInARow++
			continue
		}

		if currentRune == comma && inMatch && isInt(previousMatch) {
			numbersMatchedInARow = 0
			firstNumberMatched = true
			continue
		}

		if currentRune == closeBracket && inMatch && isInt(previousMatch) {
			fmt.Printf("%d * %d = %d\n", firstNumber, secondNumber, firstNumber*secondNumber)
			total += firstNumber * secondNumber
		}

		numbersMatchedInARow = 0
		firstNumberMatched = false
		firstNumber = 0
		secondNumber = 0
		inMatch = false
	}

	return total
}

func runeToInt(r rune) int {
	return int(r - '0')
}

func isInt(r rune) bool {
	return r >= '0' && r <= '9'
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
