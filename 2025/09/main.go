package main

import (
	"bufio"
	"container/heap"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	MAX = 65535
)

type PositionHeap []AreaInfo

func (h PositionHeap) Len() int { return len(h) }

func (h PositionHeap) Less(i, j int) bool { return h[i].area < h[j].area }

func (h PositionHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *PositionHeap) Push(x any) {
	*h = append(*h, x.(AreaInfo))
}

func (h *PositionHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *PositionHeap) Peek() any {
	x := (*h)[0]
	return x
}

type AreaInfo struct {
	area int
	i    int
	j    int
}

// -----------------------

type Edges struct {
	horizontal []Edge
	vertical   []Edge
}

type Edge struct {
	location int
	min      int
	max      int
}

func buildEdges(points Positions) *Edges {
	n := min(len(points.x), len(points.y))

	vertical := make([]Edge, n/2)
	horizontal := make([]Edge, n/2)

	for i := range n {
		j := (i + 1) % n

		x1, y1 := points.x[i], points.y[i]
		x2, y2 := points.x[j], points.y[j]

		if x1 == x2 {
			min, max := min(y1, y2), max(y1, y2)
			vertical = append(vertical, Edge{location: x1, min: min, max: max})
		} else {
			min, max := min(x1, x2), max(x1, x2)
			horizontal = append(horizontal, Edge{location: y1, min: min, max: max})
		}
	}

	sort.Slice(vertical, func(i, j int) bool { return vertical[i].location < vertical[j].location })
	sort.Slice(horizontal, func(i, j int) bool { return horizontal[i].location < horizontal[j].location })

	return &Edges{
		horizontal,
		vertical,
	}
}

func is_in_polygon(x2, y2 int, v_edges []Edge) bool {
	inside := false

	for _, e := range v_edges {
		edge_x2 := e.location * 2

		if edge_x2 <= x2 {
			continue
		}

		edge_min_y2 := e.min * 2
		edge_max_y2 := e.max * 2

		if y2 > edge_min_y2 && y2 < edge_max_y2 {
			inside = !inside
		}
	}

	return inside
}

func edges_intersect(x [2]int, y [2]int, edges Edges) bool {
	min_x, max_x := x[0], x[1]
	min_y, max_y := y[0], y[1]

	v_start := sort.Search(len(edges.vertical), func(i int) bool {
		return edges.vertical[i].location > min_x
	})
	v_end := sort.Search(len(edges.vertical), func(i int) bool {
		return edges.vertical[i].location >= max_x
	})

	for _, edge := range edges.vertical[v_start:v_end] {
		if max(edge.min, min_y) < min(edge.max, max_y) {
			return true
		}
	}

	h_start := sort.Search(len(edges.horizontal), func(i int) bool {
		return edges.horizontal[i].location > min_y
	})
	h_end := sort.Search(len(edges.horizontal), func(i int) bool {
		return edges.horizontal[i].location >= max_y
	})

	for _, edge := range edges.horizontal[h_start:h_end] {
		if max(edge.min, min_x) < min(edge.max, max_x) {
			return true
		}
	}

	return false
}

// -----------------------

type Positions struct {
	x []int
	y []int
}

func getInput(filename string) Positions {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var x_list []int
	var y_list []int

	for scanner.Scan() {
		coords := strings.Split(scanner.Text(), ",")
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])

		x_list = append(x_list, x)
		y_list = append(y_list, y)
	}
	return Positions{
		x: x_list,
		y: y_list,
	}
}

func abs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}

func partOne(filename string) int {
	list := getInput(filename)
	maxArea := 0

	n := min(len(list.x), len(list.y))

	for i := range n {
		for j := i + 1; j < n; j++ {
			x := abs(list.x[i] - list.x[j])
			y := abs(list.y[i] - list.y[j])
			a := x * y
			if a > maxArea {
				maxArea = a
			}
		}
	}

	return maxArea
}

func partTwo(filename string) int {
	list := getInput(filename)

	edges := buildEdges(list)

	n := min(len(list.x), len(list.y))

	h := &PositionHeap{}
	heap.Init(h)

	for i := range n {
		for j := i + 1; j < n; j++ {
			x := abs(list.x[i]-list.x[j]) + 1
			y := abs(list.y[i]-list.y[j]) + 1
			a := x * y

			if h.Len() < MAX {
				heap.Push(h, AreaInfo{area: a, i: i, j: j})
			} else {
				minItem := h.Peek().(AreaInfo)
				if a > minItem.area {
					heap.Pop(h)
					heap.Push(h, AreaInfo{area: a, i: i, j: j})
				}
			}
		}
	}

	var candidates []AreaInfo

	for h.Len() > 0 {
		candidates = append(candidates, heap.Pop(h).(AreaInfo))
	}

	sort.Slice(candidates, func(i, j int) bool { return candidates[i].area > candidates[j].area })

	for _, item := range candidates {
		x1 := list.x[item.i]
		x2 := list.x[item.j]

		y1 := list.y[item.i]
		y2 := list.y[item.j]

		mid_x := x1 + x2
		mid_y := y1 + y2

		if !is_in_polygon(mid_x, mid_y, edges.vertical) {
			continue
		}

		x := [2]int{min(x1, x2), max(x1, x2)}
		y := [2]int{min(y1, y2), max(y1, y2)}

		if !edges_intersect(x, y, *edges) {
			return item.area
		}
	}

	return 0
}

func main() {
	res := partTwo("input")
	println(res)
}
