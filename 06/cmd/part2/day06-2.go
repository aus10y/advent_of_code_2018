package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
)

var MaxDist = 10000

func readInput(f string) []byte {
	input, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}

	return input
}

type Coordinate struct {
	x, y int
}

func parseCoordinate(s string) Coordinate {
	var c Coordinate

	div := strings.Index(s, ",")
	c.x, _ = strconv.Atoi(s[:div])
	c.y, _ = strconv.Atoi(s[div+2:])

	return c
}

func parseCoordinates(lines []string) []Coordinate {
	cs := make([]Coordinate, 0, len(lines))
	for _, line := range lines {
		cs = append(cs, parseCoordinate(line))
	}
	return cs
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func manhattanDistance(a, b Coordinate) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func (c Coordinate) manhattanDistance(c2 Coordinate) int {
	return manhattanDistance(c, c2)
}

func totalDist(c Coordinate, cs []Coordinate) int {
	d := 0
	for _, cx := range cs {
		d = d + c.manhattanDistance(cx)
	}
	return d
}

func main() {
	input := readInput("../../input.txt")
	lines := strings.Split(string(input), "\n")
	coordinates := parseCoordinates(lines)

	// Lets make a bounding box around the given coordinates.
	// To do so, we need to find the most extreme x and y values of the coordinates.
	// (position 0,0 is considered to be at the top left)
	xL, xR, yT, yB := math.MaxInt32, math.MinInt32, math.MaxInt32, math.MinInt32
	for _, c := range coordinates {
		if c.x < xL {
			xL = c.x
		} else if c.x > xR {
			xR = c.x
		}

		if c.y < yT {
			yT = c.y
		} else if c.y > yB {
			yB = c.y
		}
	}

	// Loop over every cell within the bounding box,
	// Determine whether the cell has a total distance to all coordinates within the target.
	area := 0
	for x := xL; x < xR+1; x++ {
		for y := yT; y < yB+1; y++ {
			if dist := totalDist(Coordinate{x, y}, coordinates); dist < MaxDist {
				area++
			}
		}
	}

	fmt.Printf("The size of the target region: %d\n", area)
}
