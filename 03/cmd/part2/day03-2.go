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

type Claim struct {
	ID, x, y, width, height int
}

func parseClaim(c string) Claim {
	// ID   left, top: width x height
	// #1 @ 551,185: 21x10
	claimParts := strings.Split(c[1:], "@")
	claimID, _ := strconv.Atoi(strings.TrimSpace(claimParts[0]))
	dimensions := strings.Split(strings.TrimSpace(claimParts[1]), ":")

	cornerStr := strings.TrimSpace(dimensions[0])
	areaStr := strings.TrimSpace(dimensions[1])

	corner := strings.Split(cornerStr, ",")
	area := strings.Split(areaStr, "x")

	left, _ := strconv.Atoi(corner[0])
	top, _ := strconv.Atoi(corner[1])
	width, _ := strconv.Atoi(area[0])
	height, _ := strconv.Atoi(area[1])

	return Claim{
		ID:     claimID,
		x:      left,
		y:      top,
		width:  width,
		height: height,
	}
}

func claimAreas(claim Claim) []Area {
	areas := make([]Area, 0, 50)
	for i := claim.x; i < claim.x+claim.width; i++ {
		for j := claim.y; j < claim.y+claim.height; j++ {
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

	// Associate fabric areas to the claims that cover them.
	fabricClaims := make(map[Area][]int)

	// Track claims that are not overlapping any other claims.
	standalone := make(map[int]bool)

	// Track claims that have been found to overlap other claims.
	overlapping := make(map[int]bool)

	// For each claim...
	for _, claimStr := range strings.Split(string(input), "\n") {
		// Parse the claim, find each area of fabric that it covers.
		claim := parseClaim(claimStr)
		areas := claimAreas(claim)

		// For each area in the claim...
		for _, area := range areas {
			// Check for previous claims on this area.
			previousClaims := fabricClaims[area]

			if len(previousClaims) == 0 && !overlapping[claim.ID] {
				// If there are no previous claims, and this claim has not yet
				// overlapped another, then the claim is currently standalone.
				standalone[claim.ID] = true
			} else {
				// The claim overlaps one or more others.
				// Ensure that the current claim, and the one(s) it overlaps
				// are marked as overlapping claims and removed from the standalone map.
				delete(standalone, claim.ID)
				overlapping[claim.ID] = true
				for _, claimID := range previousClaims {
					if standalone[claimID] {
						delete(standalone, claimID)
						overlapping[claimID] = true
					}
				}
			}
			fabricClaims[area] = append(previousClaims, claim.ID)
		}
	}

	if len(standalone) != 1 {
		log.Fatalf("more than one standalone claim found!")
	}

	for ID := range standalone {
		fmt.Printf("Standalone claim ID: %d\n", ID)
		break
	}
}
