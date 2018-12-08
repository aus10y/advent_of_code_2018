package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
)

func readInput(f string) []byte {
	input, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}

	return input
}

func reactUnit(a, b byte) bool {
	return a == b+32 || b == a+32
}

func filterUnit(p []byte, u byte) []byte {
	fp := make([]byte, 0, len(p))

	for _, x := range p {
		if x == u || x == u-32 {
			continue
		}
		fp = append(fp, x)
	}

	return fp
}

func reactPolymer(p []byte) []byte {
	// reacted polymer
	rp := make([]byte, 0, len(p))

	// Set up the reacted polymer by placing the first 'unit'
	rp = append(rp, p[0])

	// Iterate over the polymer.
	for i := 1; i < len(p); i++ {
		// Compare the current polymer unit to last unit of the reacted polymer.
		if reactUnit(p[i], rp[len(rp)-1]) {
			// Drop the last unit of the reacted polymer.
			rp = rp[:len(rp)-1]

			// If rp is empty, we need to jump forward.
			if len(rp) == 0 {
				if i+1 < len(p) {
					rp = append(rp, p[i+1])
					i++
				}
			}
		} else {
			rp = append(rp, p[i])
		}
	}

	return rp
}

func main() {
	polymer := readInput("../../input.txt")

	units := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

	var problemUnit byte
	shortest := math.MaxUint32
	for _, unit := range units {
		rp := reactPolymer(filterUnit(polymer, unit))
		if len(rp) < shortest {
			shortest = len(rp)
			problemUnit = unit
		}
	}

	fmt.Printf("Removing unit '%s/%s' yields the shortest polymer, with length: %d\n", string(problemUnit), string(problemUnit-32), shortest)
}
