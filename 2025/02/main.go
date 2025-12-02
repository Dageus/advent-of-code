package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Range struct {
	FistIdx int
	LastIdx int
}

var ranges []Range

func getInput(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// only one line
	scanner.Scan()

	rangesList := strings.SplitSeq(scanner.Text(), ",")
	for r := range rangesList {
		rangeValues := strings.Split(r, "-")

		first, _ := strconv.Atoi(rangeValues[0])
		last, _ := strconv.Atoi(rangeValues[1])

		ranges = append(ranges, Range{FistIdx: first, LastIdx: last})
	}
}

func partOne(filename string) int {
	getInput(filename)

	sum := 0

	for _, invalidRange := range ranges {
		for i := invalidRange.FistIdx; i <= invalidRange.LastIdx; i++ {
			stringNum := strconv.Itoa(i)
			length := len(stringNum)

			if length%2 != 0 {
				continue
			}

			mid := length / 2
			if stringNum[:mid] == stringNum[mid:] {
				sum += i
			}
		}
	}
	return sum
}

func isInvalid(stringNum string, length int) bool {
	for chunkSize := 1; chunkSize <= length/2; chunkSize++ {

		if length%chunkSize != 0 {
			continue
		}

		base := stringNum[:chunkSize]
		valid := true

		for j := chunkSize; j < length; j += chunkSize {
			if stringNum[j:j+chunkSize] != base {
				valid = false
				break
			}
		}

		if valid {
			return true
		}
	}
	return false
}

func partTwo(filename string) int {
	getInput(filename)

	sum := 0

	for _, invalidRange := range ranges {
		for i := invalidRange.FistIdx; i <= invalidRange.LastIdx; i++ {
			stringNum := strconv.Itoa(i)
			length := len(stringNum)

			// go through each case
			if isInvalid(stringNum, length) {
				sum += i
			}
		}
	}
	return sum
}

func main() {
	res := partTwo("input")
	println(res)
}
