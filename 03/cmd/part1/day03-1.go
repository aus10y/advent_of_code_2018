package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Area struct {
	x, y int
}

func claimAreas(claim string) []Area {
	// ID   left, top: width x height
	// #1 @ 551,185: 21x10
	dimensionStr := strings.Split(claim, "@")[1]
	dimensions := strings.Split(dimensionStr, ":")

	cornerStr := strings.TrimSpace(dimensions[0])
	areaStr := strings.TrimSpace(dimensions[1])

	corner := strings.Split(cornerStr, ",")
	area := strings.Split(areaStr, "x")

	left, _ := strconv.Atoi(corner[0])
	top, _ := strconv.Atoi(corner[1])
	width, _ := strconv.Atoi(area[0])
	height, _ := strconv.Atoi(area[1])

	areas := make([]Area, 0, 50)

	for i := left; i < left+width; i++ {
		for j := top; j < top+height; j++ {
			area := Area{x: i, y: j}
			areas = append(areas, area)
		}
	}

	return areas
}

func main() {
	input, err := ioutil.ReadFile("../../input.txt")
	if err != nil {
		log.Fatal(err)
	}

	claims := strings.Split(string(input), "\n")

	single := make(map[Area]bool)
	overlapping := make(map[Area]bool)
	for _, claim := range claims {
		positions := claimAreas(claim)
		for _, position := range positions {
			if single[position] {
				overlapping[position] = true
			} else {
				single[position] = true
			}
		}
	}

	fmt.Printf("Number of overlapping areas: %d\n", len(overlapping))
}
