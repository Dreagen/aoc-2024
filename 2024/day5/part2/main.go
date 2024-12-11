package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	startTime := time.Now()

	// instructionMap, err := readInstructions("../test1.txt")
	instructionMap, err := readInstructions("../input1.txt")
	if err != nil {
		fmt.Printf("Error reading list1: %s", err.Error())
		os.Exit(1)
	}

	// updates, err := readUpdates("../test2.txt")
	updates, err := readUpdates("../input2.txt")
	if err != nil {
		fmt.Printf("Error reading list1: %s", err.Error())
		os.Exit(1)
	}

	invalidUpdates := findInvalidUpdates(updates, instructionMap)

	for _, invalidUpdate := range invalidUpdates {
		for isValidUpdate(invalidUpdate, instructionMap) == false {
			fixOrder(instructionMap, invalidUpdate)
		}
	}

	total := 0
	for _, invalidUpdate := range invalidUpdates {
		invalidUpdate = fixOrder(instructionMap, invalidUpdate)
		midPoint := len(invalidUpdate) / 2
		total += invalidUpdate[midPoint]
	}

	fmt.Println()
	fmt.Printf("execution time: %d microseconds\n", time.Since(startTime).Microseconds())
	fmt.Printf("Total: %d\n", total)
}

func findInvalidUpdates(updates [][]int, instructionMap map[int]Instruction) [][]int {
	invalidUpdates := [][]int{}
	for _, update := range updates {
		isValid := isValidUpdate(update, instructionMap)

		if isValid == false {
			invalidUpdates = append(invalidUpdates, update)
		}
	}
	return invalidUpdates
}

func isValidUpdate(update []int, instructionMap map[int]Instruction) bool {
	isValid := true
	for i, page := range update {
		if isValid == false {
			break
		}

		instr, _ := instructionMap[page]

		lowers := update[:i]
		if instr.isValidLower(lowers) == false {
			isValid = false
		}

		highers := update[i+1:]
		if instr.isValidHigher(highers) == false {
			isValid = false
		}
	}

	return isValid
}

func fixOrder(instructionMap map[int]Instruction, update []int) []int {
	for i, page := range update {
		instr := instructionMap[page]

		current := i
		for j, innerPage := range update {
			if j <= i {
				continue
			}

			if instr.isInLowers(innerPage) {
				update = swap(update, current, j)
				current = j
			}
		}
	}

	return update
}

func swap(update []int, i, j int) []int {
	valueAtI := update[i]
	valueAtJ := update[j]

	update[i] = valueAtJ
	update[j] = valueAtI

	return update
}

type Instruction struct {
	key     int
	lowers  []int
	highers []int
}

func (i *Instruction) addHigher(value int) {
	i.highers = append(i.highers, value)
}

func (i *Instruction) addLower(value int) {
	i.lowers = append(i.lowers, value)
}

func (i *Instruction) isValidLower(lowerPages []int) bool {
	for _, page := range lowerPages {
		for _, higher := range i.highers {
			if higher == page {
				return false
			}
		}
	}

	return true
}

func (i *Instruction) isInLowers(page int) bool {
	for _, lower := range i.lowers {
		if page == lower {
			return true
		}
	}

	return false
}

func (i *Instruction) isInHighers(page int) bool {
	for _, higher := range i.highers {
		if page == higher {
			return true
		}
	}

	return false
}

func (i *Instruction) isValidHigher(higherPages []int) bool {
	for _, page := range higherPages {
		for _, lower := range i.lowers {
			if lower == page {
				return false
			}
		}
	}

	return true
}

func printInstructions(instructionMap map[int]Instruction) {
	for item := range instructionMap {
		instr := instructionMap[item]
		fmt.Println()
		printInstruction(instr)
	}
}

func printInstruction(instr Instruction) {
	fmt.Printf("key: %v\n", instr.key)
	fmt.Printf("highers: %v\n", instr.highers)
	fmt.Printf("lowers: %v\n", instr.lowers)
}

func printUpdates(updates [][]int) {
	for _, update := range updates {
		fmt.Print("pages: ")
		for _, page := range update {
			fmt.Printf("%d, ", page)
		}
		fmt.Println()
	}
}

func readInstructions(fileName string) (map[int]Instruction, error) {
	instructionMap := make(map[int]Instruction)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		instructions := strings.Split(line, "|")

		key, err := strconv.Atoi(instructions[0])
		if err != nil {
			return nil, err
		}

		value, err := strconv.Atoi(instructions[1])
		if err != nil {
			return nil, err
		}

		if instr, exists := instructionMap[key]; exists {
			instr.addHigher(value)
			instructionMap[key] = instr
		} else {
			instr := Instruction{key: key}
			instr.addHigher(value)
			instructionMap[key] = instr
		}

		if instr, exists := instructionMap[value]; exists {
			instr.addLower(key)
			instructionMap[value] = instr
		} else {
			instr := Instruction{key: value}
			instr.addLower(key)
			instructionMap[value] = instr
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return instructionMap, nil
}

func removeValue(slice []int, value int) []int {
	newSlice := []int{}
	for _, v := range slice {
		if v != value {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}

func readUpdates(fileName string) ([][]int, error) {
	updates := [][]int{}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		update := []int{}
		pages := strings.Split(line, ",")
		for _, page := range pages {
			updateInt, _ := strconv.Atoi(page)
			update = append(update, updateInt)
		}

		updates = append(updates, update)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return updates, nil
}
