package main

import (
	"bufio"
	"math"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

const MaxInt = math.MaxInt

type Machine struct {
	Target   int64
	Buttons  []int64
	Joltages []int64
}

func parseLights(s string) int64 {

	mask := int64(0)

	inner_content := s[1 : len(s)-1]

	for i, char := range inner_content {
		if char == '#' {
			mask |= (1 << i)
		}
	}

	return mask
}

func parseButtons(s []string) []int64 {
	buttonMapping := make([]int64, len(s))

	for i, button := range s {
		temp := int64(0)
		button = strings.Trim(button, "()")
		parts := strings.SplitSeq(button, ",")

		for p := range parts {
			if p == "" {
				continue
			}
			temp |= (1 << (p[0] - '0'))
		}

		buttonMapping[i] = temp
	}

	return buttonMapping
}

func parseJoltage(s string) []int64 {
	vars := strings.Split(s[1:len(s)-1], ",")

	joltageMapping := make([]int64, len(vars))

	for i, joltage := range vars {
		num, _ := strconv.Atoi(joltage)
		joltageMapping[i] = int64(num)
	}

	return joltageMapping
}

func minimalPresses(machine Machine) int64 {
	target := machine.Target

	type State struct {
		mask    int64
		presses int64
	}

	queue := []State{{mask: 0, presses: 0}}

	visited := map[int]bool{0: true}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.mask == target {
			return curr.presses
		}

		for _, buttonMask := range machine.Buttons {
			nextMask := curr.mask ^ buttonMask

			if !visited[int(nextMask)] {
				visited[int(nextMask)] = true
				queue = append(queue, State{
					mask:    nextMask,
					presses: curr.presses + 1,
				})
			}
		}
	}

	return 0
}

func partOne(filename string) int64 {
	buttonPresses := int64(0)

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		vars := strings.Split(scanner.Text(), " ")
		n := len(vars)

		buttonPresses += minimalPresses(Machine{
			Target:   parseLights(vars[0]),
			Buttons:  parseButtons(vars[1 : n-1]),
			Joltages: parseJoltage(vars[n-1]),
		})
	}

	return buttonPresses
}

func smallerOrEqual(a, b []int64) bool {
	for i := 0; i < len(a); i++ {
		if a[i] > b[i] {
			return false
		}
	}
	return true
}

func equalsModulo2(a, b []int64) bool {
	for i := 0; i < len(a); i++ {
		if a[i]%2 != b[i]%2 {
			return false
		}
	}
	return true
}

func isZero(a []int64) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != 0 {
			return false
		}
	}
	return true
}

// Using Gaussian Elimination
func (m *Machine) minimalPressesWithJoltage(combinations []ButtonCombination) (int64, bool) {
	if isZero(m.Joltages) {
		return int64(0), true
	}

	res := int64(MaxInt)

	for _, comb := range combinations {
		if !smallerOrEqual(comb.joltages, m.Joltages) {
			continue
		}
		if !equalsModulo2(comb.joltages, m.Joltages) {
			continue
		}

		nextMachine := Machine{
			Target:   0,
			Buttons:  []int64{},
			Joltages: make([]int64, len(m.Joltages)),
		}
		for i := range len(m.Joltages) {
			nextMachine.Joltages[i] = (m.Joltages[i] - comb.joltages[i]) / 2
		}

		rec, ok := nextMachine.minimalPressesWithJoltage(combinations)
		if !ok {
			continue
		}

		if n := 2*rec + int64(comb.nPressedButtons); n < res {
			res = n
		}
	}
	if res < MaxInt {
		return res, true
	}
	return 0, false
}

func (m *Machine) allCombinations() []ButtonCombination {
	nButtons := len(m.Buttons)
	if nButtons == 0 {
		return []ButtonCombination{{joltages: make([]int64, len(m.Joltages)), nPressedButtons: 0}}
	}

	res := make([]ButtonCombination, 0, 1<<nButtons)
	for n := 0; n < (1 << nButtons); n++ {
		counter := make([]int64, len(m.Joltages))
		nPressedButtons := 0
		for j := 0; j < nButtons; j++ {
			if (n & (1 << j)) != 0 {
				nPressedButtons++
				mask := m.Buttons[j]
				for mask > 0 {
					idx := bits.TrailingZeros64(uint64(mask))

					counter[idx]++

					// Clear that bit so we find the next one in the next iteration
					mask &= (mask - 1)
				}
			}
		}
		res = append(res, ButtonCombination{counter, nPressedButtons})
	}
	return res
}

type ButtonCombination struct {
	joltages        []int64
	nPressedButtons int
}

func partTwo(filename string) int64 {
	buttonPresses := int64(0)

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		vars := strings.Split(scanner.Text(), " ")
		n := len(vars)

		m := Machine{
			Target:   parseLights(vars[0]),
			Buttons:  parseButtons(vars[1 : n-1]),
			Joltages: parseJoltage(vars[n-1]),
		}

		combinations := m.allCombinations()

		res, ok := m.minimalPressesWithJoltage(combinations)

		if ok {
			buttonPresses += res
		}
	}

	return buttonPresses
}

func main() {
	res := partTwo("input")
	println(res)
}
