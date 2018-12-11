package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
)

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

var EquiDistant = Coordinate{math.MaxInt32, math.MaxInt32}

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

func finiteArea(c Coordinate, cs []Coordinate) bool {
	// Determine if c has finite 'area', by checking whether it's edge adjacent
	// neighbors approach any of the other coordinates.

	// Represent infinitude in cardinal direction
	n, s, e, w := true, true, true, true

	for _, x := range cs {
		if x == c {
			continue // Skip coordinate c
		}

		n = n && manhattanDistance(Coordinate{c.x, c.y + 1}, x) < manhattanDistance(Coordinate{c.x, c.y + 2}, x)
		s = s && manhattanDistance(Coordinate{c.x, c.y - 1}, x) < manhattanDistance(Coordinate{c.x, c.y - 2}, x)
		e = e && manhattanDistance(Coordinate{c.x + 1, c.y}, x) < manhattanDistance(Coordinate{c.x + 2, c.y}, x)
		w = w && manhattanDistance(Coordinate{c.x - 1, c.y}, x) < manhattanDistance(Coordinate{c.x - 2, c.y}, x)
	}

	return !(n || s || e || w)
}

func closestCoord(cell Coordinate, cs []Coordinate, cache map[Coordinate]Coordinate) Coordinate {
	var ok bool
	var closest Coordinate
	var shortest = math.MaxInt32

	if closest, ok = cache[cell]; !ok {
		// Find the closest coordinate to the current cell
		for _, coord := range cs {
			d := cell.manhattanDistance(coord)
			if d < shortest {
				shortest = d
				closest = coord
			} else if d == shortest {
				closest = EquiDistant
			}
		}
		cache[cell] = closest
	}

	return closest
}

func coordArea(c Coordinate, cs []Coordinate, cache map[Coordinate]Coordinate) int {
	// To find the 'area' of c,
	// Expand outward from c, tallying the cells that are nearest to c,
	// Stop when an entire iteration is nearest to a cell other than c.
	area := 1
	i := 1
	for {
		found := false

		// Top and bottom perimeter
		for x := -i; x < i+1; x++ {
			for _, y := range [2]int{-i, i} {
				cell := Coordinate{c.x + x, c.y + y}
				if closest := closestCoord(cell, cs, cache); closest == c {
					area++
					found = true
				}
			}
		}

		// Left and right perimeter
		for y := -i + 1; y < i; y++ {
			for _, x := range [2]int{-i, i} {
				cell := Coordinate{c.x + x, c.y + y}
				if closest := closestCoord(cell, cs, cache); closest == c {
					area++
					found = true
				}
			}
		}

		if !found {
			break
		}

		if i == 100 {
			area = -1
			break
		}

		i++
	}

	return area
}

func main() {
	input := readInput("../../input.txt")
	lines := strings.Split(string(input), "\n")
	coordinates := parseCoordinates(lines)

	// Find coordinates that don't have infinite areas.
	finite := make([]Coordinate, 0, len(coordinates))
	for _, c := range coordinates {
		if finiteArea(c, coordinates) {
			finite = append(finite, c)
		}
	}

	// For each coordinate of finite area, compute the area.
	coordCache := make(map[Coordinate]Coordinate)
	areaMap := make(map[Coordinate]int)
	for _, c := range finite {
		areaMap[c] = coordArea(c, coordinates, coordCache)
	}

	area := 0
	var largeC Coordinate
	for c, a := range areaMap {
		if a > area {
			area = a
			largeC = c
		}
	}

	fmt.Printf("Coordinate '%v' has the largest area: %d\n", largeC, area)
}
