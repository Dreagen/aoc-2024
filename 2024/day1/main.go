package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	lines = 1000
)

func main() {
	startTime := time.Now()

	list1, list2, err := readDistances("combined.txt")
	if err != nil {
		fmt.Printf("Error reading list1: %s", err.Error())
		os.Exit(1)
	}

	if len(*list1) != len(*list2) {
		fmt.Printf("Lists don't have the same lengths, list1: %d, list2: %d", len(*list1), len(*list2))
		os.Exit(1)
	}

	slices.Sort(*list1)
	slices.Sort(*list2)

	var totalForPart1 int
	var totalForPart2 int

	var matchingCount int
	for i, item1 := range *list1 {
		item2 := (*list2)[i]

		matchingCount = 0
		for _, item2loop := range *list2 {
			if item2loop == item1 {
				matchingCount++
			}
		}

		totalForPart1 += int(math.Abs(float64(item1 - item2)))
		totalForPart2 += item1 * matchingCount
	}

	duration := time.Since(startTime)

	fmt.Printf("Total for part 1: %d\n", totalForPart1)
	fmt.Printf("Total for part 2: %d\n", totalForPart2)
	fmt.Printf("Time taken: %d microseconds\n", duration.Microseconds())
}

func readDistances(fileName string) (*[]int, *[]int, error) {
	list1 := make([]int, 0, lines)
	list2 := make([]int, 0, lines)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		nums := strings.Split(line, "   ")

		list1Item, _ := strconv.Atoi(nums[0])
		list2Item, _ := strconv.Atoi(nums[1])

		list1 = append(list1, list1Item)
		list2 = append(list2, list2Item)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return &list1, &list2, nil
}
