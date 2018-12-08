package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func readInput(f string) []byte {
	input, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}

	return input
}

func react(a, b byte) bool {
	return a == b+32 || b == a+32
}

func main() {
	polymer := readInput("../../input.txt")

	// reacted polymer
	rp := make([]byte, 0, len(polymer))

	// Set up the reacted polymer by placing the first 'unit'
	rp = append(rp, polymer[0])

	// Iterate over the polymer.
	for i := 1; i < len(polymer); i++ {
		// Compare the current polymer unit to last unit of the reacted polymer.
		if react(polymer[i], rp[len(rp)-1]) {
			// Drop the last unit of the reacted polymer.
			rp = rp[:len(rp)-1]

			// If rp is empty, we need to jump forward.
			if len(rp) == 0 {
				if i+1 < len(polymer) {
					rp = append(rp, polymer[i+1])
					i++
				}
			}
		} else {
			rp = append(rp, polymer[i])
		}
	}

	fmt.Printf("The length of the reacted polymer: %d\n", len(rp))
}
