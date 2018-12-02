package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	input, err := ioutil.ReadFile("../../input.txt")
	if err != nil {
		log.Fatal(err)
	}

	boxIDs := strings.Split(string(input), "\n")

	countDouble := 0
	countTriple := 0

	for _, boxID := range boxIDs {
		charCount := make(map[rune]int)
		for _, char := range boxID {
			charCount[char] = charCount[char] + 1
		}

		seenDouble := false
		seenTriple := false
		for _, count := range charCount {
			if count == 2 && !seenDouble {
				countDouble++
				seenDouble = true
			}
			if count == 3 && !seenTriple {
				countTriple++
				seenTriple = true
			}
		}
	}

	checksum := countDouble * countTriple

	fmt.Printf("Count Double: %d\n", countDouble)
	fmt.Printf("Count Triple: %d\n", countTriple)
	fmt.Printf("Checksum: %d\n", checksum)
}
