package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	input, err := ioutil.ReadFile("../../input-1.txt")
	if err != nil {
		log.Fatal(err)
	}

	frequenceyChanges := strings.Split(strings.TrimSpace(string(input)), "\n")

	frequency := 0
	seen := make(map[int]bool)
	for {
		for _, frequencyChange := range frequenceyChanges {
			change, err := strconv.Atoi(frequencyChange)
			if err != nil {
				log.Fatal(err)
			}

			frequency += change
			if seen[frequency] {
				fmt.Printf("The first repeating frequency: %d\n", frequency)
				return
			}

			seen[frequency] = true
		}
	}
}
