package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	lines = 1000
)

func main() {
	startTime := time.Now()
	fmt.Println("day2")

	data, err := readFile("input.txt")
	// data, err := readFile("test.txt")

	if err != nil {
		fmt.Printf("Error reading input: %s", err.Error())
		os.Exit(1)
	}

	var safeReports int
	for _, report := range *data {
		if reportIsSafe(report, true) {
			safeReports++
		}
	}

	fmt.Printf("Safe report count: %d\n", safeReports)
	endTime := time.Since(startTime)
	fmt.Printf("execution time: %d microseconds\n", endTime.Microseconds())
}

func reportIsSafe(report []int, dampenerActive bool) bool {
	previous := -99
	goingUp := false

	for i, current := range report {
		if previous == -99 {
			previous = current
			continue
		}

		if current == previous {
			if dampenerActive {
				for index := range report {
					dampenedReport := removeFromSlice(report, index)
					if reportIsSafe(dampenedReport, false) {
						return true
					}
				}
			}

			return false
		}

		if i == 1 {
			goingUp = current > previous
		}

		if current > previous {
			if goingUp == false {
				if dampenerActive {
					for index := range report {
						dampenedReport := removeFromSlice(report, index)
						if reportIsSafe(dampenedReport, false) {
							return true
						}
					}
				}

				return false
			}

			if current-previous > 3 {
				if dampenerActive {
					for index := range report {
						dampenedReport := removeFromSlice(report, index)
						if reportIsSafe(dampenedReport, false) {
							return true
						}
					}
				}

				return false
			}

			goingUp = true
		}

		if current < previous {
			if goingUp == true {
				if dampenerActive {
					for index := range report {
						dampenedReport := removeFromSlice(report, index)
						if reportIsSafe(dampenedReport, false) {
							return true
						}
					}
				}

				return false
			}

			if previous-current > 3 {
				if dampenerActive {
					for index := range report {
						dampenedReport := removeFromSlice(report, index)
						if reportIsSafe(dampenedReport, false) {
							return true
						}
					}
				}

				return false
			}

			goingUp = false
		}

		previous = current
	}

	// fmt.Printf("Report %v is safe\n", report)
	return true
}

func readFile(fileName string) (*[][]int, error) {
	var list1 [][]int
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		numsAsStrings := strings.Split(line, " ")
		nums := make([]int, len(numsAsStrings))

		for i, numAsString := range numsAsStrings {
			num, _ := strconv.Atoi(numAsString)
			nums[i] = num
		}

		list1 = append(list1, nums)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &list1, nil
}

func removeFromSlice(slice []int, index int) []int {
	newSlice := make([]int, len(slice)-1)
	copy(newSlice, slice[:index])
	copy(newSlice[index:], slice[index+1:])
	return newSlice
}
