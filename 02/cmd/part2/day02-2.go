package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"unicode/utf8"
)

func main() {
	input, err := ioutil.ReadFile("../../input.txt")
	if err != nil {
		log.Fatal(err)
	}

	boxIDs := strings.Split(string(input), "\n")
	for i, boxOne := range boxIDs[:len(boxIDs)-1] {
		for _, boxTwo := range boxIDs[i+1:] {
			diff := 0
			common := make([]rune, len(boxOne)-1)
			for j, runeOne := range boxOne {
				runeTwo, _ := utf8.DecodeRuneInString(boxTwo[j:])
				if runeTwo == utf8.RuneError {
					log.Fatal(runeTwo)
				}

				if runeOne != runeTwo {
					diff++
				} else {
					common[j-diff] = runeOne
				}

				if diff > 1 {
					break
				}
			}

			if diff == 1 {
				fmt.Printf("The common letters: %s\n", string(common))
			}
		}
	}
}
